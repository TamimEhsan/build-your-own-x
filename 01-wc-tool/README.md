# Build Your Own wc Tool

## Challenge

This challenge is to build own version of the Unix command line tool wc! It should support flags `c` `m` `w` `l` meaning output the count of bytes, chars, words and lines respectively. Unless the input has multibyte characters, the byte and char count will be the same. It should also support pipe. 

```bash
> ccwc -c test.txt
  342190 test.txt
> ccwc -l test.txt
  7145 test.txt
> ccwc -w test.txt
  58164 test.txt
> ccwc -m test.txt
  339292 test.txt
> ccwc test.txt
  7145 58164 342190 test.txt
>cat test.txt | ccwc -l
 7145
```

## Solution 

### Handling flags 
To handle flags I used the handly library `flag` of golang.

### Handling the counts
To count the varieties of queries, I looped over the bytes and increamented the counts depending on whitespaces. The corner case is consecutive whitespaces.
