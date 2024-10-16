# Build Your Own Git
This one of my most challenging task ever. Here I will build a simple git or gitlet. It will be capable of adding files to index, commiting the changes and finally pushing to github. 

## Solution

### `git init`
The git init command simply creates some directory and files in the location specified. The directories and files in `.git` are
- `.git/objects/` : This folder saves the compressed version of the files
- `.git/branches/` :
- `.git/HEAD` : 
- `.git/config` :
etc. I will add the description to these as I move forward with the tasks.

### `git hash-object`: 
This command calculates the sha1 hash of the file provided. The data that is hashed is structured like 

```
<type> <size>x\00<content>
```
The type indicates the type of the object which is one of `blob` `commit` ``. The size is the actual size of the object. The hash-object is also used in other git commands too. We talk about it later. 

### `git add`

**Index Entry:**

```
0-15 time data 
16-27 dev, ino, mode
28-29 object type(4b), unused(3b),unix permission(9b)
30-42 uid(4B),gid(4B),file size(4B)
32-bit uid
  this is stat(2) data
32-bit gid
  this is stat(2) data
32-bit file size
  This is the on-disk size from stat(2), truncated to 32-bit.
Object name for the represented object
A 16-bit 'flags' field split into (high to low bits)
1-bit assume-valid flag
1-bit extended flag (must be zero in version 2)
2-bit stage (during merge)
12-bit name length if the length is less than 0xFFF; otherwise 0xFFF
is stored in this field.
(Version 3 or later) A 16-bit field, only applicable if the
"extended flag" above is 1, split into (high to low bits).
1-bit reserved for future
1-bit skip-worktree flag (used by sparse checkout)
1-bit intent-to-add flag (used by "git add -N")
13-bit unused, must be zero
```
<!--
4+4+16*4 = 72
12+12+16*3 = 72
44 49 52 43 00 00 00 02 00 00 00 02 <67 0f 83 5a
1b 7e 5f 61 67 0f 83 5a 1b 7e 5f 61 00 00 08 02
00 f6 44 e6 00 00 81 a4 00 00 03 e8 00 00 03 e8
00 00 00 0c 3b 18 e5 12 db a7 9e 4c 83 00 dd 08
ae b3 7f 8e 72 8b 8d ad 00 05 41 2e 74 78 74 00
00 00 00 00> <67 0f 83 5a 1b 7e 5f 61 67 0f 83 5a
1b 7e 5f 61 00 00 08 02 00 f6 44 e7 00 00 81 a4
00 00 03 e8 00 00 03 e8 00 00 00 0c 3b 18 e5 12
db a7 9e 4c 83 00 dd 08 ae b3 7f 8e 72 8b 8d ad
00 05 42 2e 74 78 74 00 00 00 00 00> 27 b5 be e8
25 a3 58 15 a3 68 99 81 2a 41 da 0f 5b 5f e0 f9


44 49 52 43 00 00 00 02 00 00 00 02 <67 0f 82 e2
25 15 aa 39 67 0f 82 e2 25 15 aa 39 00 00 08 02
00 f6 44 c5 00 00 81 a4 00 00 03 e8 00 00 03 e8
00 00 00 0c 3b 18 e5 12 db a7 9e 4c 83 00 dd 08
ae b3 7f 8e 72 8b 8d ad 00 05 41 2e 74 78 74 00
00 00 00 00> <67 0f 82 e2 25 15 aa 39 67 0f 82 e2
25 15 aa 39 00 00 08 02 00 f6 44 c6 00 00 81 a4
00 00 03 e8 00 00 03 e8 00 00 00 0c 3b 18 e5 12
db a7 9e 4c 83 00 dd 08 ae b3 7f 8e 72 8b 8d ad
00 05 42 2e 74 78 74 00 00 00 00 00> 77 03 28 8c
ac 79 c2 a4 91 08 56 68 7a ce f3 68 03 c2 5e bb
-->
Field	Size (bytes)
1. ctime (8 bytes): The last time a file's metadata changed.
1. mtime (8 bytes): The last time a file's data changed.
1. dev (4 bytes): Device number.
1. ino (4 bytes): Inode number.
1. mode (4 bytes): Object type and permissions.
1. uid (4 bytes): User ID of the owner.
1. gid (4 bytes): Group ID of the owner.
1. size (4 bytes): Size of the file in bytes.
1. SHA-1 (20 bytes): SHA-1 hash of the file's contents.
1. flags (2 bytes): Flags, including the length of the file path.
-------------------
62 byte  
file name with path null terminated  
the whole index then will be padded with null byte to make multiple of 8