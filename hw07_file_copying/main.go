package main

import (
	"flag"
	"fmt"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	if from == "" {
		fmt.Println("No filepath to copy is specified. Please, use the -from flag.")
		return
	}

	if from == "" {
		fmt.Println("No destination filepath is specified. Please, use the -to flag.")
		return
	}

	err := Copy(from, to, offset, limit)
	if err != nil {
		panic(err)
	}
}
