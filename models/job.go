package models

import (
	"encoding/json"

	uuid "github.com/satori/go.uuid"
)

type Job struct {
	Id      string `json:"id"`
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func ToJob(jobJs string) (*Job, error) {
	job := &Job{}
	err := json.Unmarshal([]byte(jobJs), job)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func NewJob(payload string, jobType string) Job {
	return Job{
		Id:      uuid.NewV4().String(),
		Payload: payload,
		Type:    jobType,
	}
}

func ToJson(job Job) (string, error) {
	serialized, err := json.Marshal(job)
	if err != nil {
		return "", err
	}
	js := string(serialized[:])
	return js, nil
}
