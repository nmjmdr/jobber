package recoverer

import (
	"github.com/go-redis/redis"
	"github.com/nmjmdr/jobber/common/constants"
	"github.com/nmjmdr/jobber/common/models"
	"github.com/nmjmdr/jobber/dlock"
	"github.com/pkg/errors"
)

type Recoverer interface {
	Recover() error
}

type recoverer struct {
	client *redis.Client
}

func NewRecoverer(client *redis.Client) Recoverer {
	return &recoverer{
		client: client,
	}
}

func (r *recoverer) Recover() error {
	// read head of processing queue
	results, err := r.client.LRange(constants.InProcessQueue, 0, 0).Result()
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
	lock := dlock.NewLock(r.client.Pipeline())
	isLocked, err := lock.IsLocked(job.Id)
	if err != nil {
		return errors.Wrap(err, "Could not check for lock in recover")
	}

	if isLocked {
		return nil
	}

	// Pushing to worker queue and removing from in process queue should not cause any issues
	// There is a chance that we end up with two entries for the job in in_process_queue
	// but then we remove one of them below, the latest one is locked by the worker and we wont delete it
	_, err = r.client.RPush(constants.WorkerQueueName(job.Type), results[0]).Result()
	if err != nil {
		return errors.Wrap(err, "Unable to push job %s to worker queue")
	}
	_, err = r.client.LRem(constants.InProcessQueue, 1, results[0]).Result()
	if err != nil {
		return errors.Wrap(err, "Unable to delete job %s from perocessing queue")
	}
	return nil
}
