package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"unicode/utf8"
)

func main() {
	countBytes := flag.Bool("c", false, "count bytes")
	countLines := flag.Bool("l", false, "count lines")
	countChars := flag.Bool("m", false, "count characters")
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
		charsCount int64
	)

	buf := make([]byte, 32*1024)
	var leftOver []byte
	for {
		n, err := file.Read(buf)

		data := buf[:n]
		if len(leftOver) > 0 {
			data = append(leftOver, data...)
			leftOver = nil
		}

		bytesCount += int64(n)

		for len(data) > 0 {
			if !utf8.FullRune(data) {
				leftOver = append(leftOver, data...)
				break
			}

			r, size := utf8.DecodeRune(data)
			charsCount++

			if r == '\n' {
				linesCount++
			}

			data = data[size:]
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
	if *countChars {
		counts = append(counts, charsCount)
	}

	for _, count := range counts {
		fmt.Printf("%d ", count)
	}

	fmt.Println(fileName)
}
