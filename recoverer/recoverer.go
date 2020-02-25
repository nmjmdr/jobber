package recoverer

import (
	"log"

	"github.com/go-redis/redis"
	"github.com/nmjmdr/jobber/common/constants"
	"github.com/nmjmdr/jobber/common/models"
	"github.com/nmjmdr/jobber/dlock"
)

type Recoverer interface {
	Start()
	Stop()
}

type recoverer struct {
	client *redis.Client
	quitCh chan bool
}

func NewRecoverer(client *redis.Client) Recoverer {
	return &recoverer{
		client: client,
	}
}

func (r *recoverer) recover() {
	// read head of processing queue
	results, err := r.client.LRange(constants.InProcessQueue, 0, 0).Result()
	if err != nil && err != redis.Nil {
		log.Printf("Error while trying to recover jobs in Recoverer: %s", err)
		return
	}

	if err == redis.Nil || results == nil || len(results) == 0 {
		// nothing to recover return
		return
	}

	var job *models.Job
	job, err = models.ToJob(results[0])
	if err != nil {
		log.Printf("Error while serializing job %s, Error %s", results[0], err)
		return
	}

	// is there a active lock on the job?
	lock := dlock.NewLock(r.client.Pipeline())
	isLocked, err := lock.IsLocked(job.Id)
	if err != nil {
		log.Printf("Could not check for lock in recover, Error: %s", err)
		return
	}
	//push job to worker queue
	if isLocked {
		return
	}

	// Pushing to worker queue and removing from in process queue should not cause any issues
	// There is a chance that we end up with two entries for the job in in_process_queue
	// but then we remove one of them below, the latest one is locked by the worker and we wont delete it
	_, err = r.client.RPush(constants.WorkerQueueName(job.Type), results[0]).Result()
	if err != nil {
		log.Printf("Unable to push job %s to worker queue, Error: %s", err)
		return
	}
	_, err = r.client.LRem(constants.InProcessQueue, 1, results[0]).Result()
	if err != nil {
		log.Printf("Unable to delete job %s from perocessing queue, Error: %s", err)
		return
	}
}

func (r *recoverer) Start() {
	for {
		select {
		case _ = <-r.quitCh:
			break
		default:
			r.recover()
		}
	}
}

func (r *recoverer) Stop() {
	r.quitCh <- true
}
