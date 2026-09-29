package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	countBytes := flag.Bool("c", false, "Output number of bytes in a file")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "no file provided")
		os.Exit(1)
	}

	fileName := args[0]
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer file.Close()

	var bytesCount int64

	buf := make([]byte, 32*1024)
	for {
		n, err := file.Read(buf)
		bytesCount += int64(n)

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	if *countBytes {
		fmt.Println(bytesCount, fileName)
	}
}
