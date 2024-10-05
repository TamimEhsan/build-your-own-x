# Build Your Own Compression Tool
This challenge is to build your own command like tool to compress text files.  

## The Challenge 
In the early 1950s David Huffman developed an algorithm to find the optimal prefix code that can be used for lossless data compression.

Given there is usually an unequal distribution of character occurrences in text this can then be used to compress data by giving the most commonly occurring characters the shortest prefix.

For example if we have the string aaabbc, it would normally take up 6 bytes, but if we assign each character a variable length code, with the most frequently occurring character has the shortest code we might give them the following codes:
```
a: 1
b: 01
c: 10
```
and we could reduce the string aaabbc (six bytes) to 111010110 (nine bits). It’s not quite that simple though as we need to ensure that the codes are prefix-free, that is the bit string representing one character is not a prefix of a bit string representing another.

## Run locally
Download the test file based on Les Misérables from [here](https://www.gutenberg.org/files/135/135-0.txt).

Make sure you have go installed in your system. To encode a file, run
```
go run encode.go <input-file> <output-file>
```
To decode a file, run
```
go run decode.go <input-file> <output-file>
```

## Solution

The overall idea is to build a huffman encoding of all the input characters. Then use that to encode the file. The encoding itself needs to be added with the encoded data as a header. 
During decoding phase, the header is read first, then another huffman tree is made with that. Using that tree the file is decoded.

### [Encoding] Building a node
This be the nodes of the huffman tree. This will store the character information too.

### [Encoding] Building a binary heap
In order to use huffman encoding, we need to take the two minimum value in a set, add them and put them in the set again. We can use a linear search to find the two minimum value in O(n) time or use binary heap to do that in O(logn). Go provides an interface `container/heap`. By implementing the interface, we have our binary min heap. The binary heap will use the node struct which itself is the huffman tree. 

### [Encoding] Building the huffman tree
This will be created using the binary heap. By repeatedly poping two minimum weight node and pushing a new node created with the poped nodes as child we create the huffman tree. 

### [Encoding] The encoding itself
Create a map of the runes and their encoding. Output the encoding to the output file. Read streams of runes from the input file and output streams of the encoded runes to the output file. 
**The problem:** Go buffered io only works with bytes. As the encoding is of variable length, we can't use bytes or runes directly. Solution? create a byte as buffer, write bit by bit to the buffer. Whenever the byte is full, write that to the output file. 

### [Decoding] Building the huffman tree, again
Read the headers from the encoded file. Build a huffman tree using the header information.

### [Decoding] The decoding itself
Read from the encoded file, decode using the header and write to the decoded file. Read bit by bit and move through the tree. Whenever a leaf node is reached, output the character stored there and reposition to the root of the tree and continute.

**The problem:** Same as before, we can't read a single bit. Create a buffer as before. Read from the buffer, whenever the buffer is empty, read new bytes from the file. 
