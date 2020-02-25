package constants

func WorkerQueueName(jobType string) string {
	return jobType + "_queue"
}

const InProcessQueue = "InProcess_queue"
