package queue

import (
	"context"
	"log"
	"os"
	"strconv"
)

// Worker struct
type Worker struct {
	Name      string
	TaskQueue chan Task
}

// Task interface
type Task interface {
	Action(ctx context.Context) error
	ExecCount() int
	AddExecCount()
}

var (
	executionLimit = 2
)

func init() {
	if os.Getenv("RETRY_EXECUTION_LIMIT") != "" {
		limit, err := strconv.Atoi(os.Getenv("RETRY_EXECUTION_LIMIT"))
		if err == nil {
			executionLimit = limit
		}
	}
}

// NewWorker creates a new Worker with a channel for tasks
func NewWorker(name string) *Worker {
	return &Worker{
		Name:      name,
		TaskQueue: make(chan Task),
	}
}

func (w *Worker) StartWorker(ctx context.Context) {
	log.Printf("starting %s worker", w.Name)
	go func() {
		defer log.Printf("stopping %s worker", w.Name)
		for {
			select {
			case task := <-w.TaskQueue:
				if err := task.Action(ctx); err != nil && task.ExecCount() < executionLimit {
					task.AddExecCount()
					w.TaskQueue <- task
				}
			case <-ctx.Done():
				return // Exit the loop if the context is cancelled
			}
		}
	}()
}
