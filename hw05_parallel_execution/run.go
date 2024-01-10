package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// Проверяем, что количество горутин больше нуля
	if n <= 0 {
		return ErrErrorsLimitExceeded
	}

	tasksChan := make(chan Task)
	errorsChan := make(chan error)
	wg := sync.WaitGroup{}

	// Запускаем n рабочих горутин
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasksChan { // Обрабатываем задачи из канала
				if err := task(); err != nil {
					errorsChan <- err // Отправляем ошибку в канал, если задача завершилась с ошибкой
				}
			}
		}()
	}

	// Отдельная горутина для мониторинга ошибок
	go func() {
		errorCount := 0
		for err := range errorsChan {
			if err != nil {
				errorCount++
				if m > 0 && errorCount >= m {
					close(tasksChan) // Закрываем канал задач, если превышен лимит ошибок
					return
				}
			}
		}
	}()

	// Распределяем задачи по каналу
	for _, task := range tasks {
		tasksChan <- task
	}
	close(tasksChan)

	wg.Wait()
	close(errorsChan) // Закрываем канал ошибок после того, как все горутины завершились

	// Проверяем, был ли превышен лимит ошибок.
	if m > 0 && len(errorsChan) >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
