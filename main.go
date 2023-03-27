package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"./processor"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go_sum <path-to-log-directory>")
		os.Exit(1)
	}

	dirPath := os.Args[1]

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		os.Exit(1)
	}

	totalSize := processor.ProcessFiles(files, dirPath)

	fmt.Printf("Total size: %.2f GB\n", float64(totalSize)/1e+9)
}