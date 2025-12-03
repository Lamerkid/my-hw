package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("if maxErrorCount=0, than return error", func(t *testing.T) {
		tasksCount := 5
		tasks := make([]Task, 0, tasksCount)

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 0

		err := Run(tasks, workersCount, maxErrorsCount)
		require.Error(t, err, "Must be ErrErrorsLimitExceeded")
	})

	t.Run("taskCount < workersCount and taskCount < maxErrorCount", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				return err
			})
		}

		workersCount := 50
		maxErrorsCount := 20

		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err)
	})

	t.Run("concurrency test", func(t *testing.T) {
		tasksCount := 100

		var runningTasks int32
		var maxConcurrent int32

		tasks := make([]Task, 0, tasksCount)

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				concurrent := atomic.AddInt32(&runningTasks, 1)

				if current := atomic.LoadInt32(&maxConcurrent); concurrent > current {
					atomic.CompareAndSwapInt32(&maxConcurrent, current, concurrent)
				}

				time.Sleep(10 * time.Millisecond)
				atomic.AddInt32(&runningTasks, -1)
				return nil
			})
		}

		workersCount := 10
		Run(tasks, workersCount, 1)

		require.Eventually(t, func() bool {
			return atomic.LoadInt32(&maxConcurrent) > 1
		}, 1000*time.Millisecond, 10*time.Millisecond, "tasks should run concurrently")
	})
}
