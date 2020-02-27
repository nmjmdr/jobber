package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis"
	"github.com/nmjmdr/jobber/dispatcher"
)

type RequestBody struct {
	JobType     string      `json:"jobType"`
	PayloadMap  interface{} `json:"payload"`
	payloadJson string
}

type Response struct {
	JobId string `json:"jobId"`
}
// TODO: use HTTP middleware and http request context to do this
func parseBody(r *http.Request) (*RequestBody, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var requestBody RequestBody
	err = json.Unmarshal(b, &requestBody)
	if err != nil {
		return nil, err
	}
	payloadBytes, err := json.Marshal(requestBody.PayloadMap)
	if err != nil {
		return nil, err
	}
	requestBody.payloadJson = string(payloadBytes)
	return &requestBody, nil
}

func Handler(dispatcher dispatcher.Dispatcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := parseBody(r)
		if err != nil {
			http.Error(w, "Request body is not in the right format", http.StatusBadRequest)
			return
		}
		jobId, err := dispatcher.Post(requestBody.payloadJson, requestBody.JobType)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to post job, Error: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		response := Response{JobId: jobId}
		responseBytes, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(responseBytes)
	}
}

func main() {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", DB: 0})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Unable to ping redis error : %s\n", err.Error())
	} else {
		log.Println("> connected to redis")
	}
	handler := Handler(dispatcher.NewFifoDispatcher(client))
	log.Println("> started FIFO dispatcher")

	port := 3000
	log.Printf("> starting http server on port %d\n", port)
	r := chi.NewRouter()
	r.Post("/jobs", handler)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		log.Fatalln(err)
	}
}
