package reader_test

import (
	"bufio"
	"interview/src/generator"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

const (
	file        = "data.txt"
	alphabet    = "abcdefghijklmnopqrstuvwxyz"
	line_length = 7
	total_lines = 10_000_000
)

func TestFileCreation(test *testing.T) {
	file := uuid.New().String() + ".txt"

	// test creation time

	start := time.Now()
	generator.GenerateRandomStrings(file, []byte(alphabet), line_length, total_lines)
	duration := time.Since(start)

	test.Logf("Total time to generate strings: %s", duration)

	f, err := os.Open(file)

	if err != nil {
		test.Error("File was not created")
	}

	// go will go through the defers like a stack
	defer f.Close()
	defer os.Remove(file)

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	/*
		1. test the length of each line
		2. test the total amount of lines
	*/

	var count int = 0
	for scanner.Scan() {
		text := scanner.Text()

		if len(text) != line_length {
			test.Errorf("Unequal line length for string %s with given length of %d", text, line_length)
		}

		count++
	}

	if count != total_lines {
		test.Errorf("Expected %d lines, encountered only %d lines", total_lines, count)
	}
}
