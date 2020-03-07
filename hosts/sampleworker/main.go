package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/nmjmdr/jobber/common/redisqueue"
	"github.com/nmjmdr/jobber/dlock"
	"github.com/nmjmdr/jobber/worker"
)

const runEvery = 1 * time.Second

func main() {
	fmt.Println("Here")
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", DB: 0})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Unable to ping redis error : %s\n", err.Error())
	} else {
		log.Println("> connected to redis")
	}

	w := worker.NewWorker(
		"sample",
		20*time.Second,
		func(payload string) (string, error) {
			return payload + " - done", nil
		},
		func(result string, err error) {
			fmt.Println("Sample Worker: ", result)
		},
		redisqueue.NewRedisClientQueue(client),
		dlock.NewLock(client),
	)

	ticker := time.NewTicker(runEvery)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			select {
			case <-ticker.C:
				err := w.Work()
				if err != nil {
					log.Printf("Sample worker: Received error while trying to get jobs to work on: %s\n", err.Error())
				}
			case <-done:
				fmt.Println("\nStop signal received. Stopping sample worker")
				ticker.Stop()
				wg.Done()
				break
			}
		}
	}()
	wg.Wait()
}
