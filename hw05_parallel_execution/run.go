package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	var errorCount int32

	// Loop throgh runners
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Loop through all tasks until list is empty or maxErrors exceeded
			for {
				mu.Lock()
				if len(tasks) == 0 {
					mu.Unlock()
					return
				}
				currTask := tasks[0] // use first element
				tasks = tasks[1:]    // remove first element
				mu.Unlock()

				if err := currTask(); err != nil {
					atomic.AddInt32(&errorCount, 1)
				}
				if int(atomic.LoadInt32(&errorCount)) >= m {
					return
				}
			}
		}()
	}

	wg.Wait()

	if int(errorCount) >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
