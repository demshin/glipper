package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/atotto/clipboard"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s /path/to/your/dir", os.Args[0])
	}

	toClipPath := os.Args[1]

	fmt.Printf("Copying %s to clipboard\n", toClipPath)

	collectedContent, err := collectContentFromDir(toClipPath)

	if err != nil {
		log.Fatalf("Failed to process the directory '%s': %v", toClipPath, err)
	}

	fmt.Printf("Collected content length: %d bytes\n", len(collectedContent))
	fmt.Println("Writing to clipboard...")

	// Write to clipboard
	err = clipboard.WriteAll(collectedContent)
	if err != nil {
		log.Fatalf("Error writing to clipboard: %v\n", err)
	}
	
	fmt.Println("All files have been copied to the clipboard")
	
	// Small delay before exit
	time.Sleep(500 * time.Millisecond)
}

func collectContentFromDir(dirPath string) (string, error) {
	var allContent strings.Builder
	fileCount := 0
	
	// Add header with metadata
	allContent.WriteString("# GLIPPER OUTPUT\n")
	allContent.WriteString("# Generated: " + time.Now().Format("2006-01-02 15:04:05") + "\n")
	allContent.WriteString("# Source: " + dirPath + "\n")
	allContent.WriteString("\n")

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip hidden directories
		if info.IsDir() && shouldSkipDir(info) {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			// Skip very large files (larger than 1 MB)
			if info.Size() > 1024*1024 {
				fmt.Printf("Skipping large file: %s (%.2f MB)\n", path, float64(info.Size())/(1024*1024))
				return nil
			}
			
			// Get relative path for more compact display
			relPath, err := filepath.Rel(dirPath, path)
			if err != nil {
				relPath = path // If couldn't get relative path, use full path
			}
			
			fileContent, err := readFileContent(path)
			if err != nil {
				fmt.Printf("Warning: Failed to read file '%s': %v\n", path, err)
				return nil // Continue processing even if a file is inaccessible
			}

			// Check if the file can be considered text
			isText, _ := isTextFile(fileContent)
			if isText {
				appendFileContent(&allContent, relPath, fileContent)
				fileCount++
			} else {
				// For binary files, only add file information
				allContent.WriteString(fmt.Sprintf("## File: %s\n", relPath))
				allContent.WriteString("(Binary file, contents skipped)\n\n")
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	fmt.Printf("Total processed files: %d\n", fileCount)
	return allContent.String(), nil
}

// Checks if the content is primarily text
// Returns: (is text, content type "ascii" or "utf8")
func isTextFile(content string) (bool, string) {
	// Threshold for determining a text file (percentage of printable characters)
	const textThreshold = 0.7
	
	if len(content) == 0 {
		return true, "ascii"
	}
	
	// If the file is not valid UTF-8, check only ASCII characters
	if !utf8.ValidString(content) {
		// Check percentage of ASCII characters
		asciiCount := 0
		for i := 0; i < len(content); i++ {
			if content[i] <= 127 && (content[i] >= 32 || content[i] == '\n' || content[i] == '\r' || content[i] == '\t') {
				asciiCount++
			}
		}
		
		return float64(asciiCount)/float64(len(content)) >= textThreshold, "ascii"
	}
	
	// For UTF-8 strings, use unicode.IsPrint to check all characters
	printableCount := 0
	totalCount := 0
	
	for _, r := range content {
		totalCount++
		
		// Check if the character is printable or whitespace
		if unicode.IsPrint(r) || unicode.IsSpace(r) {
			printableCount++
		}
	}
	
	return float64(printableCount)/float64(totalCount) >= textThreshold, "utf8"
}

func shouldSkipDir(info os.FileInfo) bool {
	isHidden := info.IsDir() && len(info.Name()) > 0 && info.Name()[0] == '.'
	return isHidden
}

func readFileContent(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("Failed to read file '%s': %v", path, err)
	}
	return string(content), nil
}

func appendFileContent(builder *strings.Builder, path string, content string) {
	builder.WriteString(fmt.Sprintf("## File: %s\n", path))
	builder.WriteString("```\n")
	builder.WriteString(content)
	builder.WriteString("\n```\n\n")
}
