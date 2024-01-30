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

	tasksCh := make(chan Task, len(tasks))
	errCh := make(chan error, n)
	done := make(chan struct{})
	var wg sync.WaitGroup
	defer close(done)

	// Запуск воркеров с учетом канала done
	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(tasksCh, errCh, done, &wg)
	}

	// Отправка задач в канал задач
	for _, task := range tasks {
		tasksCh <- task
	}
	close(tasksCh)

	// Ожидание завершения воркеров, чтобы не блокировать основную горутину
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Обработка ошибок
	var errCount int
	for err := range errCh {
		if err != nil {
			errCount++
			if m > 0 && errCount >= m {
				return ErrErrorsLimitExceeded
			}
		}
	}

	return nil
}

func worker(tasksCh <-chan Task, errCh chan<- error, done <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case task, ok := <-tasksCh:
			if !ok {
				return
			}
			err := task()
			select {
			case errCh <- err:
			case <-done:
				return
			}
		case <-done:
			return
		}
	}
}
