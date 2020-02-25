package worker

import (
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/nmjmdr/jobber/common/constants"
	"github.com/nmjmdr/jobber/common/models"
	"github.com/nmjmdr/jobber/dlock"
)

type Worker interface {
	Start()
	Stop()
}

type Handle func(payload string) (string, error)
type PostResult func(result string, err error)

type worker struct {
	jobType          string
	visiblityTimeout time.Duration
	client           *redis.Client
	handle           Handle
	postResult       PostResult
	quitCh           chan bool
}

func (w *worker) work() {

	/* The steps we follow are:
	1. Read the head of the queue
	2. Try and lock the job
	3. If you cant, return, go back to waiting
	4. If locked, then RPOPLPUSH to in_process_queue
	5. This way we do not have to implement a transaction
	6. The recoverer will not be able to recover, before the worker can lock the job
	*/

	// pop from worker queue
	queue := constants.WorkerQueueName(w.jobType)

	results, err := w.client.LRange(queue, 0, 0).Result()
	if err != nil && err != redis.Nil {
		log.Printf("Error getting jobs from worker queue: %s", err)
		return
	}

	if err == redis.Nil || len(results) == 0 {
		return
	}

	var job *models.Job
	job, err = models.ToJob(results[0])
	if err != nil {
		log.Printf("Could not serialize job from queue: %s, result is: %s", err, results[0])
		return
	}

	// no jobs
	if err == redis.Nil {
		return
	}

	// we have got the job now, we should lock it, so that recoverer and other workers we are working on it
	lock := dlock.NewLock(w.client.Pipeline())
	locked, err := lock.Lock(job.Id, w.visiblityTimeout)
	if err != nil {
		log.Printf("Error encountred while trying to get for job: %s, Error: %s", job.Id, err)
		return
	}
	if !locked {
		// some other worker locked it, return
		return
	}

	// Pop and push to in_process_queue
	_, err = w.client.RPopLPush(queue, constants.InProcessQueue).Result()
	if err != nil && err != redis.Nil {
		log.Printf("Error getting jobs from worker queue: %s", err)
		return
	}

	// process job
	result, err := w.handle(job.Payload)
	w.postResult(result, err)
}

func (w *worker) Start() {
	for {
		select {
		case _ = <-w.quitCh:
			break
		default:
			w.work()
		}
	}
}

func (w *worker) Stop() {
	w.quitCh <- true
}

func NewWorker(jobType string,
	visiblityTimeout time.Duration,
	handle Handle,
	postResult PostResult,
	client *redis.Client,
) Worker {
	return &worker{
		jobType:          jobType,
		visiblityTimeout: visiblityTimeout,
		handle:           handle,
		postResult:       postResult,
		quitCh:           make(chan bool),
		client:           client,
	}
}
