package dispatcher

import (
	"github.com/nmjmdr/jobber/common/constants"
	"github.com/nmjmdr/jobber/common/models"
	"github.com/nmjmdr/jobber/common/rediswrapper"
	"github.com/pkg/errors"
)

type fifoDispatcher struct {
	client rediswrapper.Redis
}

func (f *fifoDispatcher) Post(payload string, jobType string) (string, error) {
	job := models.NewJob(payload, jobType)
	queue := constants.WorkerQueueName(jobType)

	json, err := models.ToJson(job)
	if err != nil {
		return "", err
	}

	_, err = f.client.RPush(queue, json)
	if err != nil {
		return "", errors.Wrap(err, "Unable to push job worker queue")
	}

	return job.Id, nil
}

func NewFifoDispatcher(client rediswrapper.Redis) Dispatcher {
	return &fifoDispatcher{
		client: client,
	}
}
