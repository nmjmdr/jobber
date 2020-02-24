package worker

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/nmjmdr/jobber/constants"
	"github.com/nmjmdr/jobber/dlock"
	"github.com/nmjmdr/jobber/models"
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

func (w *worker) fetch() (*models.Job, error) {
	// perform this operation in a transaction
	pipe := w.client.TxPipeline()
	// pop from worker queue
	queue := constants.WorkerQueueName(w.jobType)
	cmd := pipe.RPopLPush(queue, constants.InProcessQueue)
	result, err := cmd.Result()
	if err != nil && err != redis.Nil {
		log.Printf("Error getting jobs from worker queue: %s", err)
		return nil, err
	}
	// no jobs
	if err == redis.Nil {
		return nil, nil
	}

	var job *models.Job
	job, err = models.ToJob(result)
	if err != nil {
		log.Print(fmt.Sprintf("Could not serialize job from queue: %s, result is: ", err, result))
		return nil, err
	}

	// we have got the job now, we should lock it, so that recoverer knows we are working on it
	lock := dlock.NewLock(pipe)
	locked, err := lock.Lock(job.Id, w.visiblityTimeout)
	if err != nil {
		log.Printf("Error encountred while trying to get for job: %s, Error: %s", job.Id, err)
		return nil, err
	}

	if !locked {
		// this should never happen,
		return nil, errors.New("Job was popped by worker, but was unable to lock it")
	}
	return job, nil

}

func (w *worker) work() {

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
