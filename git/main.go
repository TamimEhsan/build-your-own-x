package main

import (
	"compress/zlib"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"syscall"
)

type indexEntry struct {
	ctimeSec  int
	ctimeNsec int
	mtimeSec  int
	mtimeNsec int
	dev       int
	ino       int
	mode      int
	uid       int
	gid       int
	size      int
	sha1      []byte
	flags     int
	path      string
}

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
	indexEntries := make([]indexEntry, 0)
	for _, filename := range files {
		indexEntry := indexEntry{}
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

		fileInfo, err := os.Stat(filename)
		if err != nil {
			log.Fatalf("Failed to get file info: %v", err)
		}
		stat := fileInfo.Sys().(*syscall.Stat_t)
		indexEntry.ctimeSec = int(stat.Ctim.Sec)
		indexEntry.ctimeNsec = int(stat.Ctim.Nsec)
		indexEntry.mtimeSec = int(stat.Mtim.Sec)
		indexEntry.mtimeNsec = int(stat.Mtim.Nsec)
		indexEntry.dev = int(stat.Dev)
		indexEntry.ino = int(stat.Ino)
		indexEntry.mode = int(stat.Mode)
		indexEntry.uid = int(stat.Uid)
		indexEntry.gid = int(stat.Gid)
		indexEntry.size = int(stat.Size)
		indexEntry.path = filename
		indexEntry.sha1 = []byte(fileHash)
		indexEntry.flags = len(filename)

		indexEntries = append(indexEntries, indexEntry)

	}
	writeToIndex(indexEntries)
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

	zlibWriter, _ := zlib.NewWriterLevel(objectFile, zlib.DefaultCompression)
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

func writeToIndex(indexEntries []indexEntry) {
	file, err := os.OpenFile(path.Join("git", "index"), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to open index file: %v", err)
	}
	defer file.Close()
	header := []byte("DIRC")
	header = append(header, paddInteger(2, 4)...)
	headerBytes := append([]byte(header), paddInteger(len(indexEntries), 4)...)
	file.Write(headerBytes)
	for _, entry := range indexEntries {
		fmt.Println(entry.path)
		file.Write(paddInteger(entry.ctimeSec, 4))
		file.Write(paddInteger(entry.ctimeNsec, 4))
		file.Write(paddInteger(entry.mtimeSec, 4))
		file.Write(paddInteger(entry.mtimeNsec, 4))
		file.Write(paddInteger(entry.dev, 4))
		file.Write(paddInteger(entry.ino, 4))
		file.Write(paddInteger(entry.mode, 4)) // probably some issue
		file.Write(paddInteger(entry.uid, 4))
		file.Write(paddInteger(entry.gid, 4))
		file.Write(paddInteger(entry.size, 4))
		file.Write(entry.sha1) // probably some issue
		file.Write(paddInteger(entry.flags, 2))
		file.Write([]byte(entry.path))
		file.Write([]byte{0})

		pad := (8 - ((len(entry.path) + 1) % 8)) % 8
		file.Write(make([]byte, pad))

	}

}

func paddInteger(n int, size int) []byte {
	if size == 2 {
		b := make([]byte, size)
		binary.BigEndian.PutUint16(b, uint16(n))
		return b
	}

	b := make([]byte, size)
	binary.BigEndian.PutUint32(b, uint32(n))
	return b
}
