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
	// This is to support executing commands to lock as part of transaction
	pipe redis.Pipeliner
}

func newLocker(pipe redis.Pipeliner) Lock {
	return &locker{
		randomValue: uuid.NewV4().String(),
		pipe:        pipe,
	}
}

func (l *locker) Lock(id string, expiry time.Duration) (bool, error) {
	isSet, err := l.pipe.SetNX(id, l.randomValue, expiry).Result()
	return isSet, err
}

func (l *locker) Unlock(id string) error {
	val, err := l.pipe.Get(id).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	// check if it is this instance of LockExp that owns the lock
	if val == l.randomValue {
		return l.pipe.Del(id).Err()
	} else {
		return errors.New("Cannot unlock, not the owner of lock")
	}
}

func (l *locker) IsLocked(id string) (bool, error) {
	result, err := l.pipe.Exists(id).Result()
	return result == 1, err
}

func NewLock(pipe redis.Pipeliner) Lock {
	return newLocker(pipe)
}
