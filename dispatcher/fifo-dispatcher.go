package dispatcher

import (
	"github.com/go-redis/redis"
	"github.com/nmjmdr/jobber/common/constants"
	"github.com/nmjmdr/jobber/common/models"
)

type fifoDispatcher struct {
	client *redis.Client
}

func (f *fifoDispatcher) Post(payload string, jobType string) (string, error) {
	job := models.NewJob(payload, jobType)
	queue := constants.WorkerQueueName(jobType)

	json, err := models.ToJson(job)
	if err != nil {
		return "", err
	}

	f.client.RPush(queue, json)

	return job.Id, nil
}

func NewFifoDispatcher(client *redis.Client) Dispatcher {
	return &fifoDispatcher{
		client: client,
	}
}
