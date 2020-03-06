package recoverer

import (
	"log"

	"github.com/go-redis/redis"
	"github.com/nmjmdr/jobber/common/constants"
	"github.com/nmjmdr/jobber/common/models"
	"github.com/nmjmdr/jobber/common/redisqueue"
	"github.com/nmjmdr/jobber/dlock"
	"github.com/pkg/errors"
)

type Recoverer interface {
	Recover() error
}

type recoverer struct {
	queue  redisqueue.Queue
	locker dlock.Lock
}

func NewRecoverer(queue redisqueue.Queue, locker dlock.Lock) Recoverer {
	return &recoverer{
		queue:  queue,
		locker: locker,
	}
}

func (r *recoverer) Recover() error {
	// read head of processing queue
	results, err := r.queue.Peek(constants.InProcessQueue)
	if err != nil && err != redis.Nil {
		return errors.Wrap(err, "Error while trying to recover jobs in Recoverer")
	}

	if err == redis.Nil || results == nil || len(results) == 0 {
		return nil
	}

	var job *models.Job
	job, err = models.ToJob(results[0])
	if err != nil {
		return errors.Wrapf(err, "Error while serializing job %s", results[0])
	}

	// is there a active lock on the job?
	isLocked, err := r.locker.IsLocked(job.Id)

	if err != nil {
		return errors.Wrap(err, "Could not check for lock in recover")
	}

	if isLocked {
		return nil
	}

	// Pushing to worker queue and removing from in process queue should not cause any issues
	// There is a chance that we end up with two entries for the job in in_process_queue
	// but then we remove one of them below, the latest one is locked by the worker and we wont delete it
	err = r.queue.Push(constants.WorkerQueueName(job.Type), results[0])
	if err != nil {
		return errors.Wrap(err, "Unable to push job %s to worker queue")
	}
	err = r.queue.Remove(constants.InProcessQueue, 1, results[0])
	if err != nil {
		return errors.Wrap(err, "Unable to delete job %s from perocessing queue")
	}
	log.Printf("Receovered job: %s\n", job.Id)
	return nil
}
