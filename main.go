package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func countWordsInFile(filename string, wg *sync.WaitGroup, results chan<- map[string]int) {
	// decrement the wait group when the function is done
	defer wg.Done()

	// open the file
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening file %s: %v\n", filename, err)
		return
	}
	// close the file when the function is done
	defer file.Close()

	// create a scanner to read the file
	scanner := bufio.NewScanner(file)
	wordCount := 0

	// read the file line by line
	for scanner.Scan() {
		// read the line
		line := scanner.Text()

		// split the line into words (with whitespace)
		words := strings.Fields(line)

		// add the number of words to the word count
		wordCount += len(words)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file %s: %v\n", filename, err)
		return
	}

	// send the result to the channel
	results <- map[string]int{filename: wordCount}
}

func main() {
	// get all files in textfiles directory
	files, err := filepath.Glob("textfiles/*.text")
	if err != nil {
		log.Fatalf("Error getting files: %v\n", err)
	}

	totalWords := 0

	// create a wait group to wait for all files to be processed
	var wg sync.WaitGroup
	wg.Add(len(files))

	// create a channel to receive the results
	results := make(chan map[string]int, len(files))

	// process the files concurrently
	for _, filename := range files {
		go countWordsInFile(filename, &wg, results)
	}

	// wait for all the files to be processed
	wg.Wait()

	// close the channel to indicate that no more results will be sent
	close(results)

	// collect and display the results
	fmt.Println("Word count results:")
	for result := range results {
		for filename, count := range result {
			fmt.Printf("%s: %d words\n", filename, count)
			totalWords += count
		}
	}

	fmt.Printf("\nTotal words in all files: %d\n", totalWords)
}
