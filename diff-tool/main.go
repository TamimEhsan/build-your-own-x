package main

import (
	"bufio"
	"crypto/sha256"
	"flag"
	"fmt"
	"os"
)

type flags struct {
	filenames []string
}

func main() {

	flags := handleFlags()

	file1, err1 := os.Open(flags.filenames[0])
	file2, err2 := os.Open(flags.filenames[1])
	if err1 != nil || err2 != nil {
		fmt.Println("Error opening files")
		os.Exit(1)
	}
	defer file1.Close()
	defer file2.Close()

	// declare two arrays to store the lines of the files
	lines1 := make([]string, 0)
	lines2 := make([]string, 0)
	readLines(file1, &lines1)
	readLines(file2, &lines2)
	file1.Seek(0, 0)
	file2.Seek(0, 0)

	// compare the lines of the files
	diff := ""
	LCS(&lines1, &lines2, &diff)
	fmt.Println(diff)
	Diff(file1, file2, diff)

}

func handleFlags() flags {

	flag.Parse()
	filenames := flag.Args()

	if len(filenames) != 2 {
		fmt.Println("Please provide two filenames to compare")
		os.Exit(1)
	}

	return flags{filenames}
}

func readLines(file *os.File, lines *[]string) {
	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			break
		}
		hash := sha256.Sum256([]byte(line))
		*lines = append(*lines, fmt.Sprintf("%x", hash))
	}
}

func Diff(file1 *os.File, file2 *os.File, diff string) {

	bufIn1, bufIn2 := bufio.NewReader(file1), bufio.NewReader(file2)
	lineCount1, lineCount2 := 1, 1

	for i := 0; i < len(diff); i++ {
		c := diff[i]
		if c == 'N' {
			_, _ = bufIn1.ReadString('\n')
			_, _ = bufIn2.ReadString('\n')
			lineCount1++
			lineCount2++
			continue
		}
		j := i + 1
		for j < len(diff) && diff[j] != 'N' {
			j++
		}
		deletions, insertions := 0, 0

		for k := i; k < j; k++ {
			if diff[k] == 'D' {
				deletions++
			} else {
				insertions++
			}
		}

		printMeta(deletions, insertions, lineCount1, lineCount2)
		printDiff(deletions, insertions, bufIn1, bufIn2, lineCount1, lineCount2)
		lineCount1 += deletions
		lineCount2 += insertions
		i = j - 1

	}
}

func printMeta(deletions, insertions, lineCount1, lineCount2 int) {
	for k := 0; k < deletions; k++ {
		fmt.Printf("%d,", lineCount1+k)
	}
	if deletions == 0 {
		fmt.Printf("%da", lineCount1-1)
	} else if insertions == 0 {
		fmt.Printf("\bd%d,", lineCount2-1)
	} else {
		fmt.Printf("\bc")
	}
	for k := 0; k < insertions; k++ {
		fmt.Printf("%d,", lineCount2+k)
	}
	fmt.Print("\b \n")
}

func printDiff(deletions, insertions int, bufIn1, bufIn2 *bufio.Reader, lineCount1, lineCount2 int) {
	for k := 0; k < deletions; k++ {
		line, _ := bufIn1.ReadString('\n')
		fmt.Printf("< %s", line)
		lineCount1++
	}
	if deletions != 0 && insertions != 0 {
		fmt.Println("---")
	}
	for k := 0; k < insertions; k++ {
		line, _ := bufIn2.ReadString('\n')
		fmt.Printf("> %s", line)
		lineCount2++
	}
}

func LCS(X, Y *[]string, diff *string) {
	m := len(*X)
	n := len(*Y)
	L := make([][]int, m+1)
	for i := 0; i <= m; i++ {
		L[i] = make([]int, n+1)
	}

	for i := 0; i <= m; i++ {
		for j := 0; j <= n; j++ {
			if i == 0 || j == 0 {
				L[i][j] = 0
			} else if (*X)[i-1] == (*Y)[j-1] {
				L[i][j] = L[i-1][j-1] + 1
			} else {
				L[i][j] = max(L[i-1][j], L[i][j-1])
			}
		}
	}

	i := m
	j := n
	for i > 0 && j > 0 {
		if (*X)[i-1] == (*Y)[j-1] {
			*diff += "N"
			i--
			j--
		} else if L[i][j-1] >= L[i-1][j] {
			*diff += "I"
			j--
		} else {
			*diff += "D"
			i--
		}
	}

	for i > 0 {
		*diff += "D"
		i--
	}

	for j > 0 {
		*diff += "I"
		j--
	}

	reverse(diff)
}

func reverse(s *string) {
	runes := []rune(*s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	*s = string(runes)
}
