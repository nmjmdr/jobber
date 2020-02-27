package rediswrapper

import "github.com/go-redis/redis"

type Redis interface {
	RPush(string, string) (int64, error)
}

type redisClientWrapper struct {
	client *redis.Client
}

func NewRedisClientWrapper(client *redis.Client) Redis {
	return &redisClientWrapper{
		client: client,
	}
}

func (r *redisClientWrapper) RPush(key string, value string) (int64, error) {
	return r.client.RPush(key, value).Result()
}
