package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func updateFile(filePath string, currentTime time.Time, counter int) error {
	fileContents, err := readLines(filePath)
	if err != nil {
		return err
	}

	updatedContents := updateContent(fileContents, currentTime, counter)

	return writeLines(filePath, updatedContents)
}

func readLines(filePath string) ([]string, error) {
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

func updateContent(contents []string, currentTime time.Time, counter int) []string {
	for i, line := range contents {
		if strings.HasPrefix(line, "Текущее время:") {
			contents[i] = fmt.Sprintf("Текущее время: %s", currentTime.Format("2006-01-02 15:04:05"))
		} else if strings.HasPrefix(line, "Счетчик:") {
			contents[i] = fmt.Sprintf("Счетчик: %d", counter)
		}
	}
	return contents
}

func writeLines(filePath string, lines []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
