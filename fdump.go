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
	"log"
	"os"
)

const usageHdr = "usage: fdump -c <count> -o <offset> -h <file>"

func main() {
	log.SetFlags(0)
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
	fmt.Println("file opened")
	fmt.Printf("type(file) = %T\n", file)	// type = *os.File
}

func checkFile(path string) {
	st, err := os.Stat(path)
	if err != nil {
		log.Fatalf("fatal: cannot stat %s", path)
	}
	if st.IsDir() {
		log.Fatalf("fatal: %s is a directory", path)
	}
}
