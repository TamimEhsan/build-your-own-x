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
