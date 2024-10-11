package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"io"
	"os"
)

type node struct {
	r      rune
	weight int
	left   *node
	right  *node
}

func (n *node) isLeaf() bool {
	return n.left == nil && n.right == nil
}

type nodeHeap []*node

func (h nodeHeap) Len() int { return len(h) }
func (h nodeHeap) Less(i, j int) bool {
	return h[i].weight < h[j].weight
}
func (h nodeHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *nodeHeap) Push(x interface{}) {
	*h = append(*h, x.(*node))
}

func (h *nodeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type huffmanTree struct {
	root *node
}

func (t *huffmanTree) buildHuffmanTree(runeCount map[rune]int) {
	nodes := &nodeHeap{}
	for r, w := range runeCount {
		*nodes = append(*nodes, &node{r: r, weight: w})
	}

	heap.Init(nodes)

	for nodes.Len() > 1 {
		node1 := heap.Pop(nodes).(*node)
		node2 := heap.Pop(nodes).(*node)

		heap.Push(nodes, &node{
			weight: node1.weight + node2.weight,
			left:   node1,
			right:  node2,
		})
	}

	t.root = nodes.Pop().(*node)
}


func encode(node *node, prefix string, encoding map[rune]string) {
	if node.isLeaf() {
		encoding[node.r] = prefix
		return
	}

	encode(node.left, prefix+"0", encoding)
	encode(node.right, prefix+"1", encoding)
}

func main() {
	inputFilePath := os.Args[1]
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		fmt.Println("File not found")
		os.Exit(1)
	}
	defer inputFile.Close()

	bufIn := bufio.NewReader(inputFile)
	runeCount := make(map[rune]int)
	for {
		rune, _, err := bufIn.ReadRune()
		if err == io.EOF {
			break
		}
		runeCount[rune]++
	}
	inputFile.Seek(0, 0)


	tree := huffmanTree{}
	tree.buildHuffmanTree(runeCount)
	// traverse(tree.root, "")
	encoding := make(map[rune]string)
	encode(tree.root, "", encoding)

	outputPath := "encoded.txt"
	out, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("Error creating file")
		os.Exit(1)
	}
	defer out.Close()
	bufOut := bufio.NewWriter(out)

	for r, code := range encoding {
		fmt.Fprintln(bufOut, r, code)
	}
	fmt.Fprintln(bufOut, 0)

	// fmt.Println("encoding")
	bufIn.Reset(inputFile)
	var bitBuffer byte // To accumulate bits
	var bitCount uint8 // To count the number of bits accumulated so far
	for {
		r, _, err := bufIn.ReadRune()
		if err == io.EOF {
			break
		}

		binaryStr := encoding[r]
		for _, char := range binaryStr {

			if char == '1' {
				bitBuffer |= (1 << (7 - bitCount)) // Set the bit at the correct position
			}
			bitCount++
			if bitCount == 8 {
				bufOut.WriteByte(bitBuffer)
				bitBuffer = 0
				bitCount = 0
			}
		}
	}
	if bitCount > 0 {
		bufOut.WriteByte(bitBuffer)
	}

	bufOut.Flush()

}



