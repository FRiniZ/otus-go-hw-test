package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func worker(ch <-chan Task, wg *sync.WaitGroup, mu *sync.Mutex, countErrors *int) {
	for {
		t, ok := <-ch
		if !ok {
			break
		}
		if t() != nil {
			mu.Lock()
			*countErrors++
			mu.Unlock()
		}
	}
	wg.Done()
}

func Run(tasks []Task, n, m int) error {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	ch := make(chan Task)
	var err error
	var countErrors int

	for w := 0; w < n; w++ {
		wg.Add(1)
		go worker(ch, wg, mu, &countErrors)
	}

	for _, t := range tasks {
		ch <- t
		if m > 0 {
			mu.Lock()
			if countErrors >= m {
				err = ErrErrorsLimitExceeded
			}
			mu.Unlock()
			if err != nil {
				break
			}
		}
	}
	close(ch)
	wg.Wait()
	return err
}
