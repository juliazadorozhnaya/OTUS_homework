package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	taskChan := make(chan Task)
	var wg sync.WaitGroup
	var errorCount int32

	for _, task := range tasks {
		taskChan <- task
	}
	close(taskChan)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(taskChan, &wg, &errorCount)
	}

	wg.Wait()

	if errorCount >= int32(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func worker(tasks <-chan Task, wg *sync.WaitGroup, counter *int32) {
	defer wg.Done()

	for task := range tasks {
		err := task()
		if err != nil {
			atomic.AddInt32(counter, 1)
		}
	}
}
