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

## Run locally
Download the test file based on The Art of War from [here](https://www.dropbox.com/scl/fi/d4zs6aoq6hr3oew2b6a9v/test.txt?rlkey=20c9d257pxd5emjjzd1gcbn03&e=1&dl=0).

Make sure you have go installed in your system. To use a file as source, run
```
go run . [-c] [-m] [-w] [-l] <input-file>
```
To use stdin as source, run
```
cat <input-file> | go run . [-c] [-m] [-w] [-l]
```

## Solution 

### Handling flags 
To handle flags the handly library `flag` of golang is used.

### Handling the counts
To count the varieties of queries, loop over the bytes and increament the counts depending on whitespaces. The corner case is consecutive whitespaces. But this creates problem to handle multibyte character. Using `rune` data type of golang makes it easier to handle this.

