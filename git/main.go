package main

import (
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

func main() {
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	hashObjectCmd := flag.NewFlagSet("hash-object", flag.ExitOnError)
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)

	// statusCmd := flag.NewFlagSet("status", flag.ExitOnError)

	// commitCmd := flag.NewFlagSet("commit", flag.ExitOnError)

	// diffCmd := flag.NewFlagSet("diff", flag.ExitOnError)

	// pushCmd := flag.NewFlagSet("push", flag.ExitOnError)

	objectType := hashObjectCmd.String("t", "blob", "The type of the object")

	switch os.Args[1] {
	case "init":
		initCmd.Parse(os.Args[2:])
		directory := initCmd.Arg(0)
		git_init(directory)
	case "hash-object":
		hashObjectCmd.Parse(os.Args[2:])
		file := hashObjectCmd.Arg(0)
		git_hash_object(file, *objectType)
	case "add":
		addCmd.Parse(os.Args[2:])
		files := addCmd.Args()
		git_add(files)
	}

}

func createDir(path string) {
	// check if the directory exists
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return
	}
	if err := os.Mkdir(path, 0755); err != nil {
		log.Fatalf("Failed to create directory %s: %v", path, err)
	}
}

func createFile(path string) {
	if _, err := os.Create(path); err != nil {
		log.Fatalf("Failed to create file %s: %v", path, err)
	}
}

func git_init(directory string) {
	// create a directory
	fmt.Println(directory)
	if directory == "" || directory == "." {
		directory = "."
	} else {
		createDir(directory)
	}

	// create a .git directory
	createDir(path.Join(directory, "git"))
	// create subdirectories
	for _, dir := range []string{"branches", "hooks", "info", "logs", "objects", "refs"} {
		createDir(path.Join(directory, "git", dir))
	}
	// create files
	for _, file := range []string{"HEAD", "config", "description", "index", "packed-refs"} {
		createFile(path.Join(directory, "git", file))
	}

}

func git_hash_object(filename string, objectType string) {

	var file *os.File
	err := error(nil)
	if filename == "" || filename == "-" {
		file, err = os.CreateTemp("", "buffered-content-")
		if err != nil {
			log.Fatalf("failed to create temporary file: %v", err)
		}
		defer os.Remove(file.Name())
		io.Copy(file, os.Stdin)
		file.Seek(0, 0)
	} else {
		file, err = os.Open(filename)
		if err != nil {
			log.Fatalf("Failed to open file %s: %v", filename, err)
		}
		defer file.Close()
	}

	sz := getFileSize(file)
	hash := hash_object(file, objectType, sz)
	fmt.Println(hash)
}

func git_add(files []string) {
	for _, filename := range files {
		// create a file
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("Failed to open file %s: %v", filename, err)
		}
		defer file.Close()
		sz := getFileSize(file)
		fileHash := hash_object(file, "blob", sz)
		file.Seek(0, 0)
		writeToObject(file, fileHash, "blob", sz)

	}
}

/**
 * hash_object is a function that takes a file and an object type and
 * returns the sha1 hash of the object in form "<type> <size>\x00\<content>""
 */
func hash_object(reader io.Reader, objectType string, sz int) string {

	header := fmt.Sprintf("%s %d\x00", objectType, sz)
	hasher := sha1.New()
	hasher.Write([]byte(header))

	io.Copy(hasher, reader)

	return hex.EncodeToString(hasher.Sum(nil))

}

/**
 * writeToObject is a function that takes a reader, a file hash, an object type and a size
 * and writes the object to the .git/objects directory
 * The object is written in the form "<type> <size>\x00\<content>"
*/

func writeToObject(reader io.Reader, fileHash string, objectType string, sz int) {

	createDir(path.Join("git", "objects", fileHash[:2]))
	objectFile, err := os.Create(path.Join("git", "objects", fileHash[:2], fileHash[2:]))
	if err != nil {
		log.Fatalf("Failed to create object file: %v", err)
	}
	defer objectFile.Close()

	zlibWriter, _ := zlib.NewWriterLevel(objectFile, zlib.BestCompression)
	defer zlibWriter.Close()

	header := fmt.Sprintf("%s %d\x00", objectType, sz)
	zlibWriter.Write([]byte(header))
	io.Copy(zlibWriter, reader)
	zlibWriter.Flush()
}

func getFileSize(file *os.File) int {
	fileStat, err := file.Stat()
	if err != nil {
		log.Fatalf("Failed to get file stat: %v", err)
	}
	return int(fileStat.Size())
}
