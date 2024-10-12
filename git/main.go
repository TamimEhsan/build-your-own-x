package main

import (
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
	case "hash_object":
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
	err := os.Mkdir(directory, 0755)
	if err != nil {
		panic(err)
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

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open file %s: %v", file, err)
		os.Exit(1)
	}
	defer file.Close()

	hash := hash_object(file, objectType)
	fmt.Println(hash)

}

func git_add(files []string) {
	for _, filename := range files {
		// create a file
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("Failed to open file %s: %v", filename, err)
			os.Exit(1)
		}
		defer file.Close()
		hash_object(file, "blob")

	}
}

/**
 * hash_object is a function that takes a file and an object type and
 * returns the sha1 hash of the object in form "<type> <size>\x00\<content>""
 */
func hash_object(file *os.File, objectType string) string {

	fileStat, err := file.Stat()
	if err != nil {
		log.Fatalf("Failed to get file stats: %v", err)
	}

	header := fmt.Sprintf("%s %d\x00", objectType, fileStat.Size())

	fileContent := make([]byte, fileStat.Size())
	if _, err := io.ReadFull(file, fileContent); err != nil {
		log.Fatalf("Failed to read file content: %v", err)
	}

	data := append([]byte(header), fileContent...)

	hasher := sha1.New()
	if _, err := hasher.Write(data); err != nil {
		log.Fatalf("Failed to write data to hasher: %v", err)
	}

	return hex.EncodeToString(hasher.Sum(nil))
}
