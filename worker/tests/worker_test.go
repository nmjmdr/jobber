/* The steps we follow are:
1. Read the head of the queue
2. Try and lock the job
3. If you cant, return, go back to waiting
4. If locked, then RPOPLPUSH to in_process_queue
5. This way we do not have to implement a transaction
6. The recoverer will not be able to recover, as the worker would haved locked the job
7. Process the job
8. Delete from in_process_queue
9. Delete the lock
*/

package tests

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/nmjmdr/jobber/common/constants"
	"github.com/nmjmdr/jobber/common/models"
	"github.com/nmjmdr/jobber/common/redisqueue/mock_redisqueue"
	"github.com/nmjmdr/jobber/dlock/mock_dlock"
	"github.com/nmjmdr/jobber/worker"
)

func Test_When_A_Job_At_The_Head_Is_Locked_By_Other_Worker_It_Just_Returns(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	payload := "{}"
	jobType := "test1"

	job := models.NewJob(payload, jobType)

	jobJs, _ := models.ToJson(job)
	results := []string{jobJs}

	queue := mock_redisqueue.NewMockQueue(ctrl)

	queue.EXPECT().
		Peek(gomock.Eq(constants.WorkerQueueName(jobType))).
		Return(results, nil)

	visiblityTimeout := 1 * time.Second

	lck := mock_dlock.NewMockLock(ctrl)
	lck.EXPECT().
		Lock(gomock.Eq(job.Id), gomock.Eq(visiblityTimeout)).
		Return(false, nil)

	wrk := worker.NewWorker(jobType, visiblityTimeout, nil, nil, queue, lck)

	err := wrk.Work()

	if err != nil {
		t.Fatal(err)
	}
}

func Test_When_A_Job_At_The_Head_Is_Available_It_Locks_The_Job_And_Then_Follows_The_Neccessary_Steps(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	payload := "{}"
	jobType := "test1"

	job := models.NewJob(payload, jobType)

	jobJs, _ := models.ToJson(job)
	results := []string{jobJs}

	queue := mock_redisqueue.NewMockQueue(ctrl)

	queue.EXPECT().
		Peek(gomock.Eq(constants.WorkerQueueName(jobType))).
		Return(results, nil)

	queue.EXPECT().
		PopPush(gomock.Eq(constants.WorkerQueueName(jobType)), gomock.Eq(constants.InProcessQueue)).
		Return(nil)

	queue.EXPECT().
		Remove(gomock.Eq(constants.InProcessQueue), int64(1), gomock.Eq(jobJs)).
		Return(nil)

	visiblityTimeout := 1 * time.Second

	lck := mock_dlock.NewMockLock(ctrl)
	lck.EXPECT().
		Lock(gomock.Eq(job.Id), gomock.Eq(visiblityTimeout)).
		Return(true, nil)

	lck.EXPECT().
		Unlock(gomock.Eq(job.Id)).
		Return(nil)

	handleCalled := false
	handle := func(payload string) (string, error) {
		handleCalled = true
		return "", nil
	}
	postResult := func(result string, err error) {

	}
	wrk := worker.NewWorker(jobType, visiblityTimeout, handle, postResult, queue, lck)

	err := wrk.Work()

	if err != nil {
		t.Fatal(err)
	}

	if !handleCalled {
		t.Fatal("Handler function should have been called")
	}
}
