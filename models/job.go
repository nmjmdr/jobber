package models

import "encoding/json"

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
