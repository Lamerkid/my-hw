package main

import "fmt"

const length = 50

type ProgressBar struct {
	current int
	done    bool
}

func (pb *ProgressBar) Start() {
	pb.current = 0
	pb.done = false
	MakeProgress(pb.current)
}

func (pb *ProgressBar) Increment(amount int) {
	pb.current += amount
	if pb.current > 100 {
		pb.current = 100
	}
	MakeProgress(pb.current)
}

func (pb *ProgressBar) Finish() {
	pb.current = 100
	MakeProgress(pb.current)
	pb.done = true
	fmt.Println()
}

func (pb *ProgressBar) Cleanup() {
	if !pb.done {
		fmt.Println()
		pb.done = true
	}
}

func MakeProgress(count int) {
	n := count * length / 100

	fmt.Printf("\r[")
	for i := 0; i < n; i++ {
		fmt.Printf("#")
	}
	for i := 0; i < length-n; i++ {
		fmt.Printf(" ")
	}
	fmt.Printf("] %d%%", count)
}
