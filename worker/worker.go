package worker

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/nmjmdr/jobber/common/constants"
	"github.com/nmjmdr/jobber/common/models"
	"github.com/nmjmdr/jobber/common/redisqueue"
	"github.com/nmjmdr/jobber/dlock"
	"github.com/pkg/errors"
)

type Worker interface {
	Work() error
}

type Handle func(payload string) (string, error)
type PostResult func(result string, err error)

type worker struct {
	jobType          string
	visiblityTimeout time.Duration
	queue            redisqueue.Queue
	handle           Handle
	postResult       PostResult
	locker           dlock.Lock
}

func (w *worker) Work() error {

	/* The steps we follow are:
	1. Read the head of the queue
	2. Try and lock the job
	3. If you cant, return, go back to waiting
	4. If locked, then RPOPLPUSH to in_process_queue
	5. This way we do not have to implement a transaction
	6. The recoverer will not be able to recover, as the worker would haved locked the job
	7. Process the job
	8. Delete from in_process_queue
	9. Delete the lock
	*/

	// pop from worker queue
	queue := constants.WorkerQueueName(w.jobType)
	fmt.Println("Queue :", queue)
	results, err := w.queue.Peek(queue)

	if err != nil && err != redis.Nil {
		return errors.Wrap(err, "Error getting jobs from worker queue")
	}

	if err == redis.Nil || len(results) == 0 {
		return nil
	}

	var job *models.Job
	job, err = models.ToJob(results[0])
	if err != nil {
		return errors.Wrap(err, "Could not serialize job from queue")
	}

	// we have got the job now, we should lock it, so that recoverer and other workers we are working on it
	locked, err := w.locker.Lock(job.Id, w.visiblityTimeout)
	if err != nil {
		return errors.Wrapf(err, "Error encountred while trying to get for job: %s", job.Id)
	}
	if !locked {
		// some other worker locked it, return
		return nil
	}

	// Pop and push to in_process_queue
	err = w.queue.PopPush(queue, constants.InProcessQueue)
	if err != nil && err != redis.Nil {
		return errors.Wrap(err, "Error getting jobs from worker queue")
	}

	// process job
	result, err := w.handle(job.Payload)
	w.postResult(result, err)

	// need to delete from processing queue
	// need to then delete lock
	return nil
}

func NewWorker(jobType string,
	visiblityTimeout time.Duration,
	handle Handle,
	postResult PostResult,
	queue redisqueue.Queue,
	locker dlock.Lock,
) Worker {
	return &worker{
		jobType:          jobType,
		visiblityTimeout: visiblityTimeout,
		handle:           handle,
		postResult:       postResult,
		queue:            queue,
		locker:           locker,
	}
}
