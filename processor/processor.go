package processor

import (
	"bufio"
	"fmt"
	"time"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"
)

const numWorkers = 10

// Treats every files and sends back the total weight
func ProcessFiles(files []os.FileInfo, dirPath string) int64 {
	fileChan := make(chan string)
	sizeChan := make(chan int64)

	StartWorkerPool(fileChan, sizeChan)

	go func() {
		for _, fileInfo := range files {
			if !fileInfo.IsDir() {
				filePath := filepath.Join(dirPath, fileInfo.Name())
				fileChan <- filePath
			}
		}
		close(fileChan)
	}()

	var totalSize int64

	for size := range sizeChan {
		totalSize += size
	}

	return totalSize
}

// Starts worker pool to treat files in parallel
func StartWorkerPool(filePathsChan chan string, fileSizesChan chan int64) {
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for filePath := range filePathsChan {
				fmt.Printf("🏁 Processing %s\n", filePath)
				startTime := time.Now()
				fileSize := ProcessFile(filePath)
				elapsedTime := time.Since(startTime)
				fmt.Printf("✅ File %s processed in %s\n", filePath, elapsedTime)
				fileSizesChan <- fileSize
			}
		}()
	}

	go func() {
		wg.Wait()
		close(fileSizesChan)
	}()
}

// Treats an individual file and returns its weight
func ProcessFile(filePath string) int64 {
	logLines, err := ReadLines(filePath)
	
	if err != nil {
		fmt.Println("Error reading file:", err)
		return 0
	}

	sizePattern := regexp.MustCompile(`HTTP/1.1" 200 (\d+) `)
	var fileSize int64

	for _, line := range logLines {
		matches := sizePattern.FindStringSubmatch(line)
		if len(matches) >= 2 {
			size, err := strconv.ParseInt(matches[1], 10, 64)
			if err == nil {
				fileSize += size
			}
		}
	}

	return fileSize
}

// Reads every line of a file and returns them as an array of strings
func ReadLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}
