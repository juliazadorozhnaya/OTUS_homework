package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return ErrErrorsLimitExceeded
	}

	var mt sync.Mutex
	var errCount int
	var wg sync.WaitGroup
	tasksCh := make(chan Task)

	for _, task := range tasks {
		tasksCh <- task
	}
	close(tasksCh)

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for task := range tasksCh {
				err := task()
				mt.Lock()
				if err != nil {
					errCount++
				}
				if errCount == m {
					mt.Unlock()
					return
				}
				mt.Unlock()
			}
		}()
	}
	wg.Wait()
	if errCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
