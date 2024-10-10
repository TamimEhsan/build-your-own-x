# Build Your Own diff Tool
The diff tool has been part of every software developers tool box since it was initially released in June 1974. How’s that for legacy code that’s still providing value today!

## The Challenge
The diff tool takes inputs and outputs the difference between them line by line. 

## Solution

### LCS and Edit Graph
The edit graph is the operations needed to convert one string to another. These operations are `delete` `insert` `no change`. By deleting some lines, inserting some new lines or doing nothing we can convert a file into another file. We take the hash of each lines of the files and store them in two arrays and find the edit graph between them.

### Print diff
After we create the edit graph we print the diff file according to the output of `diff` tool