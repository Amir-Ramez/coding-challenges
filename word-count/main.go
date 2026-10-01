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

	var input io.Reader

	args := flag.Args()
	if len(args) < 1 {
		input = os.Stdin
	} else {
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer file.Close()

		input = file
	}

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
		n, err := input.Read(buf)

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

	if len(counts) == 0 {
		counts = append(counts, linesCount, wordsCount, bytesCount)
	}

	for _, count := range counts {
		fmt.Printf("%d ", count)
	}

	if input != os.Stdin {
		fmt.Println(args[0])
	}
}
