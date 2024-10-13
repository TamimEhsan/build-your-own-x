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