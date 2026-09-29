package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	countBytes := flag.Bool("c", false, "count bytes")
	countLines := flag.Bool("l", false, "count lines")
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

	var (
		bytesCount int64
		linesCount int64
	)

	buf := make([]byte, 32*1024)
	for {
		n, err := file.Read(buf)
		bytesCount += int64(n)

		for i := range n {
			if buf[i] == '\n' {
				linesCount++
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	counts := make([]int64, 0)
	if *countBytes {
		counts = append(counts, bytesCount)
	}
	if *countLines {
		counts = append(counts, linesCount)
	}

	for _, count := range counts {
		fmt.Printf("%d ", count)
	}

	fmt.Println(fileName)
}
