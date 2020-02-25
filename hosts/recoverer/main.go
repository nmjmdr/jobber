package main

import (
	"log"

	"github.com/go-redis/redis"
)

func main() {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", DB: 0})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Unable to ping redis error : %s\n", err.Error())
	} else {
		log.Println("> connected to redis")
	}
}
