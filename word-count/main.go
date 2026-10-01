package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"
)

func main() {
	countBytes := flag.Bool("c", false, "count bytes")
	countLines := flag.Bool("l", false, "count lines")
	countChars := flag.Bool("m", false, "count characters")
	countWords := flag.Bool("w", false, "count words")
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
		wordsCount int64
	)

	buf := make([]byte, 32*1024)

	var (
		leftover []byte
		inWord   = false
	)
	for {
		n, err := file.Read(buf)

		data := buf[:n]
		if len(leftover) > 0 {
			data = append(leftover, data...)
			leftover = nil
		}

		bytesCount += int64(n)

		for len(data) > 0 {
			if !utf8.FullRune(data) {
				leftover = append(leftover, data...)
				break
			}

			r, size := utf8.DecodeRune(data)
			charsCount++

			if r == '\n' {
				linesCount++
			}

			if unicode.IsSpace(r) {
				inWord = false
			} else if !inWord {
				inWord = true
				wordsCount++
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
	if *countWords {
		counts = append(counts, wordsCount)
	}

	for _, count := range counts {
		fmt.Printf("%d ", count)
	}

	fmt.Println(fileName)
}
