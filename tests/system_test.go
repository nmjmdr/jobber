package tests

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nmjmdr/jobber/common/constants"
	"github.com/nmjmdr/jobber/common/redisqueue/mock_redisqueue"
	"github.com/nmjmdr/jobber/dispatcher"
)

func Test_When_Posted_Job_Should_Be_Dispacthed_To_The_Right_Worker_Queue(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	payload := "{}"
	jobType := "test1"

	queue := mock_redisqueue.NewMockQueue(ctrl)
	queue.EXPECT().
		Push(gomock.Eq(constants.WorkerQueueName(jobType)), gomock.Any()).
		Return(nil)

	dispatcher := dispatcher.NewFifoDispatcher(queue)

	_, err := dispatcher.Post(payload, jobType)
	if err != nil {
		t.Fatal(err)
	}

}
