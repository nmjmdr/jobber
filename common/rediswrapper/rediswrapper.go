package rediswrapper

import "github.com/go-redis/redis"

type Redis interface {
	RPush(string, string) (int64, error)
	LRange(string, int64, int64) ([]string, error)
	LRem(string, int64, interface{}) (int64, error)
}

type redisClientWrapper struct {
	client *redis.Client
}

func NewRedisClientWrapper(client *redis.Client) Redis {
	return &redisClientWrapper{
		client: client,
	}
}

func (r *redisClientWrapper) LRem(key string, count int64, value interface{}) (int64, error) {
	return r.client.LRem(key, count, value).Result()
}

func (r *redisClientWrapper) LRange(key string, start int64, stop int64) ([]string, error) {
	return r.client.LRange(key, start, stop).Result()
}

func (r *redisClientWrapper) RPush(key string, value string) (int64, error) {
	return r.client.RPush(key, value).Result()
}
