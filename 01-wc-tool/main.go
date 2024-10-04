package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

type metrics struct {
	byteCount int
	charCount int
	wordCount int
	lineCount int
}
type flags struct {
	countByte bool
	countChar bool
	countWord bool
	countLine bool
}

func main() {
	flags := handleFlags()

	in, filename := getInputReader()
	if filename != "" {
		defer in.(*os.File).Close()
	}

	metrics := countMetrics(in)

	printResults(flags, metrics, filename)
}

func handleFlags() flags {
	countByte := flag.Bool("c", false, "Count bytes")
	countChar := flag.Bool("m", false, "Count characters")
	countWord := flag.Bool("w", false, "Count words")
	countLine := flag.Bool("l", false, "Count lines")

	flag.Parse()

	if !*countByte && !*countChar && !*countWord && !*countLine {
		*countByte, *countWord, *countLine = true, true, true
	}

	return flags{*countByte, *countChar, *countWord, *countLine}
}

func getInputReader() (io.Reader, string) {
	filename := ""
	if filename = flag.Arg(0); filename != "" {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return file, filename
	}
	return os.Stdin, filename
}

func isWhiteSpace(b rune) bool {
	return b == ' ' || b == '\n' || b == '\t' || b == '\r'
}

func countMetrics(in io.Reader) metrics {
	buf := bufio.NewReader(in)
	metrics := metrics{0, 0, 0, 0}
	lastByte := rune(0)
	for {
		curByte, sz, err := buf.ReadRune()
		if err == io.EOF {
			if !isWhiteSpace(lastByte) {
				metrics.wordCount++
			}
			break
		}

		metrics.byteCount += sz
		metrics.charCount++

		if isWhiteSpace(curByte) && !isWhiteSpace(lastByte) {
			metrics.wordCount++
		}
		if curByte == '\n' {
			metrics.lineCount++
		}
		lastByte = curByte
	}
	return metrics
}

func printResults(flags flags, metrics metrics, filename string) {
	if flags.countLine {
		fmt.Print(metrics.lineCount, " ")
	}

	if flags.countChar {
		fmt.Print(metrics.charCount, " ")
	}

	if flags.countWord {
		fmt.Print(metrics.wordCount, " ")
	}

	if flags.countByte {
		fmt.Print(metrics.byteCount, " ")
	}

	fmt.Println(filename)
}
