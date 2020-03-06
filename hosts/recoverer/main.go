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
	"github.com/nmjmdr/jobber/recoverer"
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

	r := recoverer.NewRecoverer(redisqueue.NewRedisClientQueue(client), dlock.NewLock(client.Pipeline()))

	ticker := time.NewTicker(runEvery)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			select {
			case <-ticker.C:
				err := r.Recover()
				if err != nil {
					log.Printf("Received error while trying to receover jobs: %s\n", err.Error())
				}
			case <-done:
				fmt.Println("\nStop signal received. Stopping recoverer")
				ticker.Stop()
				wg.Done()
				break
			}
		}
	}()
	wg.Wait()

}
