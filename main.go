package main

import (
	"fmt"
	"golang.design/x/clipboard"
	"log"
	"os"
	"path/filepath"
	"strings"
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

	if err := clipboard.Init(); err != nil {
		log.Fatalf("Failed to initialize clipboard: %v", err)
	}

	_ = clipboard.Write(clipboard.FmtText, []byte(collectedContent))

	fmt.Println("All files have been copied to the clipboard")
}

func collectContentFromDir(dirPath string) (string, error) {
	var allContent strings.Builder

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if shouldSkipDir(info) {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			fileContent, err := readFileContent(path)
			if err != nil {
				return err
			}

			if len(fileContent) > 0 {
				appendFileContent(&allContent, path, fileContent)
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return allContent.String(), nil
}

func shouldSkipDir(info os.FileInfo) bool {
	return info.IsDir() && info.Name()[0] == '.'
}

func readFileContent(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("Failed to read file '%s': %v", path, err)
	}
	return string(content), nil
}

func appendFileContent(builder *strings.Builder, path string, content string) {
	builder.WriteString(fmt.Sprintf("Path: %s\n", path))
	builder.WriteString(content)
	builder.WriteString("\n---------------------------\n\n")
}
