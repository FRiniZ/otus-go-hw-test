package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	wg := sync.WaitGroup{}
	var countErrors int
	mu := &sync.Mutex{}

	for i := 0; i < len(tasks); i += n {
		for t := 0; t < n && (i+t) < len(tasks); t++ {
			wg.Add(1)
			go func(task Task, mu *sync.Mutex) {
				defer wg.Done()
				err := task()
				if err != nil {
					mu.Lock()
					countErrors++
					mu.Unlock()
				}
			}(tasks[i+t], mu)
		}
		wg.Wait()
		if countErrors >= m {
			return ErrErrorsLimitExceeded
		}
	}

	return nil
}
