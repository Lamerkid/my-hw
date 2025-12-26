/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No directory path to read environment variables is specified.")
		return
	}

	// cmd.Execute()
	env, _ := ReadDir(os.Args[1])
	fmt.Println(env)
}
