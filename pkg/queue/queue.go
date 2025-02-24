package queue

import (
	"context"
	"log"
	"os"
	"strconv"
)

type Worker struct {
	Name      string
	TaskQueue chan Task
}

type Task interface {
	Action(ctx context.Context) error
	ExecCount() int
	AddExecCount()
	Log()
}

var (
	ExecutionLimit = 2
)

func init() {
	if os.Getenv("RETRY_EXECUTION_LIMIT") != "" {
		limit, err := strconv.Atoi(os.Getenv("RETRY_EXECUTION_LIMIT"))
		if err == nil {
			ExecutionLimit = limit
		}
	}
}

func NewWorker(name string) *Worker {
	w := Worker{
		Name:      name,
		TaskQueue: make(chan Task),
	}
	return &w
}

func (w *Worker) StartWorker(ctx context.Context) {
	log.Printf("starting %s worker", w.Name)
	for {
		select {
		case task := <-w.TaskQueue:
			if task.ExecCount() >= ExecutionLimit {
				task.Log()
				return
			}
			if err := task.Action(ctx); err != nil {
				task.AddExecCount()
				w.TaskQueue <- task
			}
		case <-ctx.Done():
			defer log.Printf("stopping %s worker", w.Name)
			return
		}
	}
}
