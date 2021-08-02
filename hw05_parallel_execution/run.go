package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	taskChan := make(chan Task)
	errorCh := make(chan error)
	doneCh := make(chan struct{})

	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(doneCh, taskChan, errorCh, &wg)
	}

	errorCount := 0
	currentTask := 0

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case err := <-errorCh:
				if err != nil {
					errorCount++
				}
			case taskChan <- tasks[currentTask]:
				currentTask++
			}
			if errorCount >= m || currentTask == len(tasks) {
				close(taskChan)
				close(doneCh)
				return
			}
		}
	}()

	wg.Wait()

	if errorCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func worker(done <-chan struct{}, tasks <-chan Task, errors chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		select {
		case <-done:
			return
		case errors <- task():
		}
	}
}
