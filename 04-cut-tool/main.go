package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type fields []int

type flags struct {
	fields    fields
	delimiter string
	filepath  string
}

func (f *fields) String() string {
	return "fields"
}

func (f *fields) Set(value string) error {
	valInt, err := strconv.Atoi(value)
	if err != nil {
		fmt.Println("flag value must be an integer")
		os.Exit(1)
	}
	*f = append(*f, valInt)
	return nil
}

func parseFlags() flags {
	fields := fields{}
	flag.Var(&fields, "f", "select only these fields")
	delimiter := flag.String("d", "\t", "use DELIM instead of TAB for field delimiter")

	flag.Parse()
	filepath := flag.Arg(0)

	return flags{fields, *delimiter, filepath}
}

func getInputReader(filepath string) io.Reader {

	if filepath == "" || filepath == "-" {
		return os.Stdin
	} else {
		file, err := os.Open(filepath)
		if err != nil {
			fmt.Println("Error opening file")
			os.Exit(1)
		}
		return file
	}

}

func cutFields(flags flags, in io.Reader) {
	bufIn := bufio.NewReader(in)
	for {
		line, err := bufIn.ReadString('\n')

		fields := strings.Split(line, flags.delimiter)
		printed := false
		for _, field := range flags.fields {
			if field <= len(fields) {
				if printed {
					fmt.Print(flags.delimiter)
				}
				fmt.Print(fields[field-1])
				printed = true
			}
		}
		fmt.Println()
		if err == io.EOF {
			break
		}

	}
}

func main() {
	flags := parseFlags()

	in := getInputReader(flags.filepath)
	if in != os.Stdin {
		defer in.(*os.File).Close()
	}

	cutFields(flags, in)
}
