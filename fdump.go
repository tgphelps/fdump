// Package fdump dumps files in the traditional hex/ASCII format.
/*
Usage: fdump -c <count> -o <offset> -h path
	-c <count>    Number of bytes to dump
	-o <offset>   Byte offset at which to start
	-h            Hex dump only
	path          File to dump
*/
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"tgphelps.com/hdump"
)

const usageHdr = "usage: fdump -c <count> -o <offset> -h <file>"

// buffSize is the number of bytes to read in one chunk.
const buffSize = 16 * 64 //must be a multiple of 16

// maxCount is the maximum number of bytes to dump. This is set to
// a number that we don't expect to ever reach.
const maxCount = 1024 * 1024 * 1024

func main() {
	log.SetFlags(0) // No date/time in messages
	var pCount = flag.Int("c", maxCount, "max bytes to dump")
	var pOffset = flag.Int("o", 0, "file offset of start")
	var pHexOnly = flag.Bool("h", false, "hex-only dump")
	flag.Parse()
	paths := flag.Args()
	if len(paths) != 1 {
		usage()
		os.Exit(1)
	}
	dump(paths[0], *pCount, *pOffset, *pHexOnly)
}

func usage() {
	fmt.Println(usageHdr)
	flag.PrintDefaults()
}

// dump performs the file dump.
func dump(path string, count int, offset int, hexOnly bool) {
	// fmt.Printf("dumping file %s offset %d count %d\n", path, offset, count)
	checkFile(path)
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// fmt.Println("file opened")
	dest := hdump.NewHdumper(os.Stdout)
	if hexOnly {
		dest.SetHexOnly(true)
	}
	dumpBytes(file, count, offset, dest)
}

// dumpBytes dumps the open file to stdout.
func dumpBytes(file *os.File, count int, offset int, dest *hdump.Hdumper) {
	// fmt.Println("dumping data from", file, "to", dest)
	buff := make([]byte, buffSize)
	if offset > 0 {
		dest.SetOffset(offset)
		file.Seek(int64(offset), io.SeekStart)
	}
	for count > 0 {
		num, err := file.Read(buff)
		if err != nil {
			if err == io.EOF {
				if num != 0 {
					log.Fatal("XXX: EOF with num > 0")
				}
				break
			} else {
				log.Fatal("error reading file:", err)
			}
		}
		if num > count {
			num = count
		}
		err = dest.DumpBytes(num, buff)
		if err != nil {
			log.Fatal("error writing dump:", err)
		}
		count -= num
	}
}

// checkFile verifies that path exists and is not a directory.
func checkFile(path string) {
	st, err := os.Stat(path)
	if err != nil {
		log.Fatalf("fatal: cannot stat %s", path)
	}
	if st.IsDir() {
		log.Fatalf("fatal: %s is a directory", path)
	}
}
