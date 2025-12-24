package main

import (
	"fmt"
	"strings"
)

const length = 50

type ProgressBar struct {
	capacity int
	current  int
	done     bool
}

func (pb *ProgressBar) Start(capacity int) {
	pb.current = 0
	pb.done = false
	pb.capacity = capacity
	MakeProgress(pb.current, pb.capacity)
}

func (pb *ProgressBar) Increment(amount int) {
	pb.current += amount
	if pb.current > pb.capacity {
		pb.current = pb.capacity
	}
	MakeProgress(pb.current, pb.capacity)
}

func (pb *ProgressBar) Finish() {
	pb.done = true
	fmt.Println()
}

func (pb *ProgressBar) Cleanup() {
	if !pb.done {
		fmt.Println()
		pb.done = true
	}
}

func MakeProgress(count, capacity int) {
	if capacity == 0 {
		fmt.Printf("\r[%s] %d/%d\tbytes 100%%", strings.Repeat(" ", length), count, capacity)
		return
	}

	percentage := count * 100 / capacity
	n := percentage * length / 100

	bar := strings.Repeat("#", n) + strings.Repeat(" ", length-n)

	fmt.Printf("\r[%s] %d/%d\tbytes %d%% ", bar, count, capacity, percentage)
}
