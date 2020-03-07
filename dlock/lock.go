package dlock

import (
	"time"

	"github.com/go-redis/redis"
)

type Lock interface {
	Lock(id string, expiry time.Duration) (bool, error)
	Unlock(id string) error
	IsLocked(id string) (bool, error)
}

type locker struct {
	randomValue string
	// This is to support executing commands to lock as part of transaction
	client *redis.Client
}

func newLocker(client *redis.Client) Lock {
	return &locker{
		client:        client,
	}
}

func (l *locker) Lock(id string, expiry time.Duration) (bool, error) {
	isSet, err := l.client.SetNX(id, "true", expiry).Result()
	return isSet, err
}

func (l *locker) Unlock(id string) error {
	return l.client.Del(id).Err()
}

func (l *locker) IsLocked(id string) (bool, error) {
	result, err := l.client.Exists(id).Result()
	return result == 1, err
}

func NewLock(client *redis.Client) Lock {
	return newLocker(client)
}