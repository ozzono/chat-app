package queue

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

// Define an Action interface
type Action interface {
	Execute(ctx context.Context) error
}

// DefaultAction is the default implementation of the Action interface
type DefaultAction struct{}

func (da DefaultAction) Execute(ctx context.Context) error {
	return nil
}

// mockTask now holds a custom action
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

func (t *mockTask) Action(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err() // Return the context's error, which will be context.Canceled
	default:
		if ca, ok := ctx.Value(customActionKey{}).(Action); ok {
			return ca.Execute(ctx)
		}
		t.AddExecCount()
		return nil
	}
}

// // Modify the Action method to call the custom action if it's set in the context
// func (t *mockTask) Action(ctx context.Context) error {
// 	if ca, ok := ctx.Value(customActionKey{}).(Action); ok {
// 		return ca.Execute(ctx)
// 	}
// 	t.AddExecCount()
// 	return nil
// }

// Define a context key for the custom action
type customActionKey struct{}

// customAction is a struct that implements the Action interface
type customAction struct {
	error error
}

// Execute implements the Action interface
func (ca *customAction) Execute(ctx context.Context) error {
	return ca.error
}

// TestWorkerHappyPath tests the happy path of task execution
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
	cancel() // Cancel the context immediately
	var task Task = &mockTask{m: &sync.Mutex{}}
	if err := task.Action(ctx); err != context.Canceled {
		t.Log("Expected context cancellation error")
		t.FailNow()
	}
}

func TestWorkerWithMultipleTasks(t *testing.T) {
	// Create a worker
	worker := NewWorker("testWorker")
	worker.StartWorker(context.Background())

	// Define a function to simulate task execution
	executeTask := func(task Task) {
		worker.TaskQueue <- task
	}

	// Create a wait group to wait for all tasks to complete
	var wg sync.WaitGroup
	wg.Add(3) // Assuming we have 3 tasks

	// Define tasks
	task1 := &mockTask{m: &sync.Mutex{}}
	task2 := &mockTask{m: &sync.Mutex{}}
	task3 := &mockTask{m: &sync.Mutex{}}

	// Execute tasks
	go func() {
		executeTask(task1)
		wg.Done()
	}()
	go func() {
		executeTask(task2)
		wg.Done()
	}()
	go func() {
		executeTask(task3)
		wg.Done()
	}()

	// Wait for all tasks to complete
	wg.Wait()

	// Verify task execution counts
	if task1.ExecCount() != 1 || task2.ExecCount() != 1 || task3.ExecCount() != 1 {
		t.Log("Tasks were not executed as expected")
		t.FailNow()
	}
}

func TestWorkerErrorHandling(t *testing.T) {
	// Create a worker
	taskWithError := &mockTask{m: &sync.Mutex{}}
	customAction := &customAction{error: fmt.Errorf("task error")}
	worker := NewWorker("testWorker")
	go worker.StartWorker(context.WithValue(context.Background(), customActionKey{}, customAction))

	// Define a task that returns an error

	// Execute the task
	worker.TaskQueue <- taskWithError

	// Wait for a short period to allow the task to be processed
	// Adjust the duration based on your worker's processing time
	time.Sleep(1 * time.Second)

	// Verify that the task was executed
	if taskWithError.ExecCount() != 1 {
		t.Log("Task was not executed")
		t.FailNow()
	}

	// Additional checks to ensure the worker handled the error appropriately
	// This could involve checking logs, metrics, or other side effects of error handling
}
