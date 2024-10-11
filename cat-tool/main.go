package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

type flags struct {
	printLineNumber         bool
	printLineNumberNonEmpty bool
	filenames               []string
}

func main() {
	flags := handleFlags()

	linenumber := 1
	for _, filename := range flags.filenames {
		in := getInputReader(filename)
		if in != os.Stdin {
			defer in.(*os.File).Close()
		}

		cat(in, flags, &linenumber)
	}

}

func handleFlags() flags {

	printLineNumber := flag.Bool("n", false, "Print line numbers")
	printLineNumberNonEmpty := flag.Bool("b", false, "Print line numbers for non-empty lines")

	flag.Parse()
	// read trailing arguments as filenames
	filenames := flag.Args()

	return flags{*printLineNumber, *printLineNumberNonEmpty, filenames}
}

func getInputReader(filename string) io.Reader {
	if filename != "" && filename != "-" {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return file
	}
	return os.Stdin
}

func cat(in io.Reader, flags flags, lineNumber *int) {

	buf := bufio.NewReader(in)
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}
		if flags.printLineNumber {
			// make the line number 6 digits wide
			fmt.Printf("%6d  %s", *lineNumber, line)
			*lineNumber++
		} else if flags.printLineNumberNonEmpty {
			if line != "\n" {
				fmt.Print(*lineNumber, "  ", line)
				*lineNumber++
			} else {
				fmt.Print(line)
			}
		} else {
			fmt.Print(line)
		}

	}

}
