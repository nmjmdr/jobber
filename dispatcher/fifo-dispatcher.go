package dispatcher

import (
	"github.com/go-redis/redis"
)

type fifoDispatcher struct {
	client *redis.Client
}

func (f *fifoDispatcher) Post(payload string, jobType string) (string, error) {
	return "", nil
}

func NewFifoDispatcher(client *redis.Client) Dispatcher {
	return &fifoDispatcher{
		client: client,
	}
}
