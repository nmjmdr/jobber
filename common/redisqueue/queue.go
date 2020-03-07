package redisqueue

import (
	"fmt"

	"github.com/go-redis/redis"
)

type Queue interface {
	Push(key string, val string) error
	Peek(key string) ([]string, error)
	Remove(key string, count int64, val interface{}) error
	PopPush(source string, destination string) error
}

type redisClientQueue struct {
	client *redis.Client
}

func NewRedisClientQueue(client *redis.Client) Queue {
	return &redisClientQueue{
		client: client,
	}
}

func (r *redisClientQueue) Remove(key string, count int64, value interface{}) error {
	v, err := r.client.LRem(key, count, value).Result()
	if v != count {
		return fmt.Errorf("Unable to delete all the keys. Requested delete of %d, but deleted: %d", count, v)
	}
	return err
}

func (r *redisClientQueue) Peek(key string) ([]string, error) {
	return r.client.LRange(key, 0, 0).Result()
}

func (r *redisClientQueue) Push(key string, value string) error {
	_, err := r.client.RPush(key, value).Result()
	return err
}

func (r *redisClientQueue) PopPush(source string, destination string) error {
	_, err := r.client.RPopLPush(source, destination).Result()
	return err
}
