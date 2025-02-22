package queue

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

type Action interface {
	Execute(ctx context.Context) error
}

type DefaultAction struct{}

func (da DefaultAction) Execute(ctx context.Context) error {
	return nil
}

type mockTask struct {
	execCount int
	m         *sync.Mutex
}

func (t *mockTask) AddExecCount() {
	t.m.Lock()
	defer t.m.Unlock()
	t.execCount++
}

func (t mockTask) ExecCount() int {
	return t.execCount
}

func (t mockTask) Log() {
	log.Println("mockTask")
}

func (t *mockTask) Action(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if ca, ok := ctx.Value(customActionKey{}).(Action); ok {
			return ca.Execute(ctx)
		}
		t.AddExecCount()
		return nil
	}
}

type customActionKey struct{}

type customAction struct {
	error error
}

func (ca *customAction) Execute(ctx context.Context) error {
	return ca.error
}

func TestWorkerHappyPath(t *testing.T) {
	ctx := context.Background()
	var task Task = &mockTask{m: &sync.Mutex{}}
	if err := task.Action(ctx); err != nil {
		t.Log("Error executing task:", err)
		t.FailNow()
	}
	if task.ExecCount() != 1 {
		t.Log("Execution count mismatch")
		t.FailNow()
	}
}

func TestContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var task Task = &mockTask{m: &sync.Mutex{}}
	if err := task.Action(ctx); err != context.Canceled {
		t.Log("Expected context cancellation error")
		t.FailNow()
	}
}

func TestWorkerErrorHandling(t *testing.T) {

	taskWithError := &mockTask{m: &sync.Mutex{}}
	customAction := &customAction{error: fmt.Errorf("task error")}
	worker := NewWorker("testWorker")
	go worker.StartWorker(context.WithValue(context.Background(), customActionKey{}, customAction))

	worker.TaskQueue <- taskWithError

	time.Sleep(1 * time.Second)

	if taskWithError.ExecCount() != 1 {
		t.Log("Task was not executed")
		t.FailNow()
	}

}
