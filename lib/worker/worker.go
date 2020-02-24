package worker

import (
	"time"

	"../dlock"

	"github.com/go-redis/redis"
)

type Worker interface {
	Start()
	Stop()
}

type Handle func(payload string) (string, error)
type PostResult func(result string, err error)

type worker struct {
	jobType          string
	visiblityTimeout time.Duration
	client           *redis.Client
	handle           Handle
	postResult       PostResult
	quitCh           chan bool
	lock             dlock.Lock
}

func (w *worker) work() {

}

func (w *worker) Start() {
	for {
		select {
		case _ = <-w.quitCh:
			break
		default:
			w.work()
		}
	}
}

func (w *worker) Stop() {

}

func NewWorker(jobType string,
	visiblityTimeout time.Duration,
	handle Handle,
	postResult PostResult,
) Worker {
	return &worker{
		jobType:          jobType,
		visiblityTimeout: visiblityTimeout,
		handle:           handle,
		postResult:       postResult,
		quitCh:           make(chan bool),
	}
}
