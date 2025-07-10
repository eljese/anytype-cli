package daemon

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// Task is a background function that runs until the given context is canceled.
type Task func(ctx context.Context) error

// taskInfo holds information about a running task
type taskInfo struct {
	cancel context.CancelFunc
	done   chan struct{}
}

// TaskManager tracks running background tasks.
type TaskManager struct {
	mu    sync.Mutex
	tasks map[string]*taskInfo
}

// defaultTaskManager is the singleton instance.
var defaultTaskManager = NewTaskManager()

// NewTaskManager returns a new task manager.
func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks: make(map[string]*taskInfo),
	}
}

// StartTask starts a new task with a unique ID.
// It returns an error if a task with that ID is already running.
func (tm *TaskManager) StartTask(id string, task Task) error {
	tm.mu.Lock()
	if _, exists := tm.tasks[id]; exists {
		tm.mu.Unlock()
		return errors.New("task already running")
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	tm.tasks[id] = &taskInfo{
		cancel: cancel,
		done:   done,
	}
	tm.mu.Unlock()

	go func() {
		defer close(done)
		if err := task(ctx); err != nil {
			fmt.Printf("Task %s exited with error: %v\n", id, err)
		}
		tm.mu.Lock()
		delete(tm.tasks, id)
		tm.mu.Unlock()
	}()
	return nil
}

// StopTask cancels a running task by its ID.
func (tm *TaskManager) StopTask(id string) error {
	tm.mu.Lock()
	task, exists := tm.tasks[id]
	tm.mu.Unlock()
	if !exists {
		return errors.New("task not found")
	}
	task.cancel()
	return nil
}

// StopAll stops every running task and waits for them to complete.
func (tm *TaskManager) StopAll() {
	tm.mu.Lock()
	var waitList []*taskInfo
	for _, task := range tm.tasks {
		task.cancel()
		waitList = append(waitList, task)
	}
	tm.mu.Unlock()

	for _, task := range waitList {
		<-task.done
	}
}

// GetTaskManager returns the singleton instance.
func GetTaskManager() *TaskManager {
	return defaultTaskManager
}
