package tests

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nmjmdr/jobber/common/constants"
	"github.com/nmjmdr/jobber/common/models"
	"github.com/nmjmdr/jobber/common/redisqueue/mock_redisqueue"
	"github.com/nmjmdr/jobber/dlock/mock_dlock"
	"github.com/nmjmdr/jobber/recoverer"
)

func Test_When_A_Job_Is_In_Process_And_Has_Active_Lock_It_Does_Not_Recover_The_Job(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	payload := "{}"
	jobType := "test1"

	job := models.NewJob(payload, jobType)

	jobJs, _ := models.ToJson(job)
	results := []string{jobJs}

	queue := mock_redisqueue.NewMockQueue(ctrl)
	// we only expect peek to be called on queue
	queue.EXPECT().
		Peek(gomock.Eq(constants.InProcessQueue)).
		Return(results, nil)

	lck := mock_dlock.NewMockLock(ctrl)
	lck.EXPECT().
		IsLocked(gomock.Eq(job.Id)).
		Return(true, nil)

	rec := recoverer.NewRecoverer(queue, lck)

	err := rec.Recover()

	if err != nil {
		t.Fatal(err)
	}

}

func Test_When_No_Job_Is_In_Process_Then_Recoverer_Just_Returns(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	results := []string{}

	queue := mock_redisqueue.NewMockQueue(ctrl)
	// we only expect peek to be called on queue
	queue.EXPECT().
		Peek(gomock.Eq(constants.InProcessQueue)).
		Return(results, nil)

	lck := mock_dlock.NewMockLock(ctrl)

	rec := recoverer.NewRecoverer(queue, lck)

	err := rec.Recover()

	if err != nil {
		t.Fatal(err)
	}

}

func Test_When_A_Job_Is_In_Process_And_Has_No_Active_Lock_It_Recovers_The_Job(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	payload := "{}"
	jobType := "test1"

	job := models.NewJob(payload, jobType)

	jobJs, _ := models.ToJson(job)
	results := []string{jobJs}

	queue := mock_redisqueue.NewMockQueue(ctrl)

	queue.EXPECT().
		Peek(gomock.Eq(constants.InProcessQueue)).
		Return(results, nil)

	queue.EXPECT().
		Push(gomock.Eq(constants.WorkerQueueName(jobType)), gomock.Eq(results[0])).
		Return(nil)

	queue.EXPECT().
		Remove(gomock.Eq(constants.InProcessQueue), gomock.Eq(int64(1)), gomock.Eq(results[0])).
		Return(nil)

	lck := mock_dlock.NewMockLock(ctrl)
	lck.EXPECT().
		IsLocked(gomock.Eq(job.Id)).
		Return(false, nil)

	rec := recoverer.NewRecoverer(queue, lck)

	err := rec.Recover()

	if err != nil {
		t.Fatal(err)
	}
}
