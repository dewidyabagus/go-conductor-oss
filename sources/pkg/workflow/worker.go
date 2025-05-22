package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/conductor-sdk/conductor-go/sdk/model"
)

type workers struct {
	taskRunner     *TaskRunner
	taskNames      []string
	inProgressTask int64
}

type WorkerDefinition struct {
	TaskName     string
	Handler      func(context.Context, *model.Task) (*model.TaskResult, error)
	BatchSize    int
	PollInterval time.Duration
	Domain       string
}

// Blocking Function
func (w *workers) RunWorkers(workerDefs []WorkerDefinition) (err error) {
	if w.taskRunner == nil {
		return errors.New("task runner not initialized")
	}
	defer func() {
		if err == nil {
			return
		}
		for _, taskName := range w.taskNames {
			w.taskRunner.Shutdown(taskName)
		}
		w.taskNames = []string{}
	}()

	for _, def := range workerDefs {
		if def.Domain == "" {
			err = w.taskRunner.StartWorker(def.TaskName, w.wrapHandler(def.Handler), def.BatchSize, def.PollInterval)

		} else {
			err = w.taskRunner.StartWorkerWithDomain(def.TaskName, w.wrapHandler(def.Handler), def.BatchSize, def.PollInterval, def.Domain)
		}
		if err != nil {
			return fmt.Errorf("start worker: %w", err)
		}
		w.taskNames = append(w.taskNames, def.TaskName)
	}

	w.taskRunner.WaitWorkers()
	return nil
}

// Note In this section, you can develop custom pre-process and post-process like adding traceId, logs, etc.
func (w *workers) wrapHandler(fn func(context.Context, *model.Task) (*model.TaskResult, error)) func(t *model.Task) (any, error) {
	return func(t *model.Task) (any, error) {
		atomic.AddInt64(&w.inProgressTask, 1)
		defer func() { atomic.AddInt64(&w.inProgressTask, -1) }()

		// For simulation only, not for implementation in production.
		raw, _ := json.Marshal(t.InputData)
		fmt.Printf("[INFO] Workflow ID: %s Task Name: %s Payload: %s\n", t.WorkflowInstanceId, t.TaskDefName, string(raw))

		ctx := context.Background()
		/////////////////////////////////////////////////////////////

		return fn(ctx, t)
	}
}

func (w *workers) Close() {
	if w.taskRunner == nil {
		return
	}

	for _, taskName := range w.taskNames {
		w.taskRunner.SetBatchSize(taskName, 0)
	}

	for atomic.LoadInt64(&w.inProgressTask) > 0 {
	}
	time.Sleep(5 * time.Second)

	for _, taskName := range w.taskNames {
		w.taskRunner.Shutdown(taskName)

		log.Println("Shutdown task:", taskName)
	}
}
