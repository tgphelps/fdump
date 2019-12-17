/*
Usage: fdump -c <count> -o <offset> -h path
	-c <count>		Number of bytes to dump
	-o <offset>		Byte offset to start at
	-h				Hex dump only. No ASCII.package fdump
	path			File to dump
*/

package main

import (
	"flag"
	"fmt"
	"hdump"
	"io"
	"log"
	"os"
)

const usageHdr = "usage: fdump -c <count> -o <offset> -h <file>"
const buffSize = 16 * 4 //must be a multiple of 16

// var buffArray [buffSize]byte
// var buff = buffArray[:]

func main() {
	log.SetFlags(0) // No date/time in messages
	var pCount = flag.Int("c", 0, "bytes to dump")
	var pOffset = flag.Int("o", 0, "file offset of start")
	var pHexOnly = flag.Bool("h", false, "hex-only dump")
	flag.Parse()
	paths := flag.Args()
	if len(paths) != 1 {
		usage()
		os.Exit(1)
	}
	// fmt.Println("args =", flag.Args())

	dump(paths[0], *pCount, *pOffset, *pHexOnly)
}

func usage() {
	fmt.Println(usageHdr)
	flag.PrintDefaults()
}

func dump(path string, count int, offset int, hexOnly bool) {
	fmt.Printf("dumping file %s offset %d count %d\n", path, count, offset)
	checkFile(path)
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// fmt.Println("file opened")
	dest := hdump.NewHdumper(os.Stdout)
	dumpBytes(file, dest)
}

func dumpBytes(file *os.File, dest *hdump.Hdumper) {
	// fmt.Println("dumping data from", file, "to", dest)
	buff := make([]byte, buffSize)
	for {
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
		err = dest.DumpBytes(num, buff)
		if err != nil {
			log.Fatal("error writing dump:", err)
		}
	}
}

func checkFile(path string) {
	st, err := os.Stat(path)
	if err != nil {
		log.Fatalf("fatal: cannot stat %s", path)
	}
	if st.IsDir() {
		log.Fatal("fatal: %s is a directory", path)
	}
}
