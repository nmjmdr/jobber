package dispatcher

import (
	"github.com/nmjmdr/jobber/common/constants"
	"github.com/nmjmdr/jobber/common/models"
	"github.com/nmjmdr/jobber/common/redisqueue"
	"github.com/pkg/errors"
)

type fifoDispatcher struct {
	queue redisqueue.Queue
}

func (f *fifoDispatcher) Post(payload string, jobType string) (string, error) {
	job := models.NewJob(payload, jobType)
	queue := constants.WorkerQueueName(jobType)

	json, err := models.ToJson(job)
	if err != nil {
		return "", err
	}

	err = f.queue.Push(queue, json)
	if err != nil {
		return "", errors.Wrap(err, "Unable to push job worker queue")
	}

	return job.Id, nil
}

func NewFifoDispatcher(queue redisqueue.Queue) Dispatcher {
	return &fifoDispatcher{
		queue: queue,
	}
}
