package util

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	once1 sync.Once
	e     *TaskExecutor
)

func getDefaultTimeout() time.Duration {
	envSize := os.Getenv("TASK_DEFAULT_TIMEOUT")
	if envSize == "" {
		return 30 * time.Second
	}
	size, err := strconv.Atoi(envSize)
	if err != nil {
		return 30 * time.Second
	}
	return time.Second * time.Duration(size)
}

func getTaskSizeFromEnv() int {
	envSize := os.Getenv("TASK_SIZE")
	if envSize == "" {
		return 1024 // 默认缓冲区大小为1024
	}
	size, err := strconv.Atoi(envSize)
	if err != nil {
		return 1024 // 如果转换失败，使用默认值
	}
	return size
}

type TaskFunc func(notify chan struct{})

type TaskExecutor struct {
	tasks     map[string]struct{}
	mu        *sync.Mutex
	queueSize int
	queue     chan struct{}
	timeout   time.Duration
}

func (e *TaskExecutor) RunTask(ctx context.Context, name string, f func() error) (chan struct{}, error) {

	//debug info , queueSize,timeout
	fmt.Println("executor spce: ", "queueSize: ", e.queueSize, "timeout: ", e.timeout)

	e.mu.Lock()
	if _, ok := e.tasks[name]; ok {
		e.mu.Unlock()
		return nil, fmt.Errorf("unrepeatable task: %s", name)
	} else {
		e.tasks[name] = struct{}{}
	}
	e.mu.Unlock()

	var cancel context.CancelFunc
	timeoutCtx := ctx

	if _, ok := ctx.Deadline(); !ok {
		timeoutCtx, cancel = context.WithTimeout(ctx, e.timeout)
		defer cancel()
	}

	// before run task
	fmt.Println("submit task ", name)
	select {
	case <-e.queue:
		fmt.Println("task: ", name, "get ticket")
	case <-timeoutCtx.Done():
		e.mu.Lock()
		delete(e.tasks, name)
		e.mu.Unlock()
		finished := make(chan struct{})
		return finished, errors.New("task has overtime")
	}

	// fmt.Println("begin run task")
	notify := make(chan struct{})
	finished := make(chan struct{})
	go func(name string) {
		defer close(notify)
		defer func() {
			e.queue <- struct{}{}
		}()
		n := 1
		ticker := time.NewTicker(time.Second * 2)
		defer ticker.Stop()
		for range ticker.C {
			err := f()
			if err != nil {
				// todo replace it with log
				fmt.Println("call task func err: ", err.Error())
				n = n * 2
				ticker.Reset(time.Duration(n) * time.Second)
			} else {
				return
			}
		}
	}(name)

	go func() {
		<-notify
		e.mu.Lock()
		delete(e.tasks, name)
		e.mu.Unlock()
		fmt.Printf("task %s has exited\n", name)
		close(finished)
	}()

	return finished, nil
}

func GetTaskExecutor() *TaskExecutor {
	once1.Do(func() {
		queueSize := getTaskSizeFromEnv()
		e = &TaskExecutor{
			tasks:     make(map[string]struct{}),
			mu:        &sync.Mutex{},
			queue:     make(chan struct{}, queueSize),
			timeout:   time.Duration(getDefaultTimeout()),
			queueSize: queueSize,
		}

		for i := 0; i < queueSize; i++ {
			e.queue <- struct{}{}
		}
		// e.queue <- struct{}{}

	})
	return e
}
