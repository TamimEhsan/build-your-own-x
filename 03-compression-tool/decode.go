package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type node struct {
	r      int
	weight int
	left   *node
	right  *node
}

func (n *node) isLeaf() bool {
	return n == nil || (n.left == nil && n.right == nil)
}

func insertNode(root *node, r int, code string) {
	for _, c := range code {
		if c == '0' {
			if root.left == nil {
				root.left = &node{}
			}
			root = root.left
		} else {
			if root.right == nil {
				root.right = &node{}
			}
			root = root.right
		}
	}
	root.r = r
}

func traverse(root *node, code string) {
	if root.isLeaf() {
		fmt.Printf("%c %s\n", root.r, code)
		return
	}
	if root.left != nil {
		traverse(root.left, code+"0")
	}
	if root.right != nil {
		traverse(root.right, code+"1")
	}
}

func main() {
	encodedFilename := os.Args[1]

	encodedFile, err := os.Open(encodedFilename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	bufIn := bufio.NewReader(encodedFile)
	root := &node{}
	count := 0
	for {
		b := 0
		fmt.Fscanf(bufIn, "%d", &b)

		if b == 0 {
			break
		}

		balsal := ""
		fmt.Fscanf(bufIn, "%s\n", &balsal)
		insertNode(root, b, balsal)
		count++
	}
	fmt.Fscanf(bufIn, "\n")

	// traverse(root, "")

	var bitBuffer string // To accumulate bits
	// To count the number of bits accumulated so far
	outputFilename := "decoded.txt"
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer outputFile.Close()
	bufOut := bufio.NewWriter(outputFile)
	for {
		b, err := bufIn.ReadByte()
		if err == io.EOF {
			break
		}

		bitBuffer += fmt.Sprintf("%08b", b)
		curNode := root
		bitPos := 0
		for itr := 0; itr < len(bitBuffer); itr++ {

			if bitBuffer[itr] == '0' {
				curNode = curNode.left
			} else {
				curNode = curNode.right
			}
			if curNode == nil {
				bufOut.Flush()
				os.Exit(1)
			}

			if curNode.isLeaf() {

				fmt.Fprint(bufOut, string(curNode.r))
				curNode = root
				bitPos = itr + 1
			}
		}

		bitBuffer = bitBuffer[bitPos:]
	}
	bufOut.Flush()
}
