package dlock

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
)

type Lock interface {
	Lock(id string, expiry time.Duration) (bool, error)
	Unlock(id string) error
	IsLocked(id string) (bool, error)
}

type locker struct {
	randomValue string
	client      *redis.Client
}

func newLocker(client *redis.Client) Lock {
	return &locker{
		randomValue: uuid.NewV4().String(),
		client:      client,
	}
}

func (l *locker) Lock(id string, expiry time.Duration) (bool, error) {
	isSet, err := l.client.SetNX(id, l.randomValue, expiry).Result()
	return isSet, err
}

func (l *locker) Unlock(id string) error {
	val, err := l.client.Get(id).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	// check if it is this instance of LockExp that owns the lock
	if val == l.randomValue {
		return l.client.Del(id).Err()
	} else {
		return errors.New("Cannot unlock, not the owner of lock")
	}
}

func (l *locker) IsLocked(id string) (bool, error) {
	_, err := l.client.Get(id).Result()
	if err == redis.Nil {
		return false, nil
	}
	return true, err
}

func NewLock(client *redis.Client) Lock {
	return newLocker(client)
}
