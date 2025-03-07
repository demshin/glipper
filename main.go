package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/atotto/clipboard"
)

// Config structure
type Config struct {
	MaxClipboardSize int  // Maximum clipboard size in bytes
	SkipBinaryFiles  bool // Skip binary files
	SkipHiddenDirs   bool // Skip hidden directories
}

// Default values
var defaultConfig = Config{
	MaxClipboardSize: 64000, // approximately 64 KB
	SkipBinaryFiles:  true,
	SkipHiddenDirs:   true,
}

// Active configuration
var config Config

func main() {
	// Load configuration
	loadConfig()

	// Parse command line flags
	var showHelp bool
	flag.IntVar(&config.MaxClipboardSize, "size", config.MaxClipboardSize, "Maximum clipboard size in bytes")
	flag.BoolVar(&config.SkipBinaryFiles, "skip-binary", config.SkipBinaryFiles, "Skip binary files")
	flag.BoolVar(&config.SkipHiddenDirs, "skip-hidden", config.SkipHiddenDirs, "Skip hidden directories")
	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.Parse()

	// If help flag is specified, show help and exit
	if showHelp {
		fmt.Println("Glipper - utility for copying file contents to clipboard")
		fmt.Printf("Usage: %s [options] [directory_path]\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Save updated configuration
	saveConfig()

	// Determine the path to copy content from
	var toClipPath string
	args := flag.Args()

	if len(args) < 1 {
		// If path is not specified, use current directory
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current directory: %v", err)
		}
		toClipPath = currentDir
		fmt.Println("Path not specified. Using current directory.")
	} else {
		toClipPath = args[0]
	}

	fmt.Printf("Copying %s to clipboard\n", toClipPath)

	collectedContent, err := collectContentFromDir(toClipPath)

	if err != nil {
		log.Fatalf("Error processing directory '%s': %v", toClipPath, err)
	}

	// Check collected content size
	contentSize := len(collectedContent)
	fmt.Printf("Collected content size: %d bytes\n", contentSize)

	// If size exceeds the limit, truncate content
	if contentSize > config.MaxClipboardSize {
		// Add message about truncation
		truncationMessage := fmt.Sprintf("\n\n# WARNING: Content has been truncated due to size limit (%d KB).\n", config.MaxClipboardSize/1024)

		// Truncate content and add truncation message
		// Leave space for truncation message
		availableSize := config.MaxClipboardSize - len(truncationMessage)
		collectedContent = collectedContent[:availableSize] + truncationMessage

		fmt.Printf("Content exceeded maximum size and was truncated to %d KB\n", config.MaxClipboardSize/1024)
	}

	fmt.Println("Writing to clipboard...")

	// Write to clipboard
	err = clipboard.WriteAll(collectedContent)
	if err != nil {
		log.Fatalf("Error writing to clipboard: %v\n", err)
	}

	fmt.Println("All files have been copied to clipboard")

	// Small delay before exit
	time.Sleep(500 * time.Millisecond)
}

// Get path to configuration file
func getConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".glipper.conf" // File in current directory if we can't get home directory
	}

	// Create the ~/.config/glipper directory if it doesn't exist
	configDir := filepath.Join(homeDir, ".config", "glipper")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Printf("Warning: Could not create config directory: %v\n", err)
		return filepath.Join(homeDir, ".glipper.conf") // Fallback to home directory
	}

	return filepath.Join(configDir, ".glipper.conf")
}

// Load configuration from file
func loadConfig() {
	// Initialize configuration with default values
	config = defaultConfig

	configPath := getConfigPath()

	// Check if configuration file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create file with default configuration
		saveConfig()
		fmt.Printf("Created default configuration file at: %s\n", configPath)
		return
	}

	// Read configuration file
	file, err := os.Open(configPath)
	if err != nil {
		fmt.Printf("Error loading config file: %v. Using default settings.\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip comments and empty lines
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		// Split line to key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Invalid configuration line
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Process different settings
		switch key {
		case "max_clipboard_size":
			if size, err := strconv.Atoi(value); err == nil {
				config.MaxClipboardSize = size
			}
		case "skip_binary_files":
			if value == "true" {
				config.SkipBinaryFiles = true
			} else if value == "false" {
				config.SkipBinaryFiles = false
			}
		case "skip_hidden_dirs":
			if value == "true" {
				config.SkipHiddenDirs = true
			} else if value == "false" {
				config.SkipHiddenDirs = false
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading config file: %v. Using partially loaded settings.\n", err)
	}
}

// Save configuration to file
func saveConfig() {
	configPath := getConfigPath()

	// Open file for writing
	file, err := os.Create(configPath)
	if err != nil {
		fmt.Printf("Error writing config file: %v\n", err)
		return
	}
	defer file.Close()

	// Write comment and settings
	fmt.Fprintln(file, "# Glipper configuration file")
	fmt.Fprintln(file, "# Generated:", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintln(file, "# Format: key=value")
	fmt.Fprintln(file, "")
	fmt.Fprintf(file, "max_clipboard_size=%d\n", config.MaxClipboardSize)
	fmt.Fprintf(file, "skip_binary_files=%v\n", config.SkipBinaryFiles)
	fmt.Fprintf(file, "skip_hidden_dirs=%v\n", config.SkipHiddenDirs)
}

func collectContentFromDir(dirPath string) (string, error) {
	var allContent strings.Builder
	fileCount := 0
	currentSize := 0
	maxSizeReached := false

	// Add header with metadata
	header := fmt.Sprintf("# GLIPPER OUTPUT\n# Generated: %s\n# Source: %s\n\n",
		time.Now().Format("2006-01-02 15:04:05"), dirPath)
	allContent.WriteString(header)
	currentSize += len(header)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// If maximum size is already reached, skip processing
		if maxSizeReached {
			return nil
		}

		// Skip hidden directories if specified in configuration
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
				// Create temporary string for file
				fileEntry := fmt.Sprintf("## File: %s\n```\n%s\n```\n\n", relPath, fileContent)
				fileEntrySize := len(fileEntry)

				// Check if adding this file will exceed the limit
				if currentSize+fileEntrySize > config.MaxClipboardSize {
					// If it will, stop processing
					maxSizeReached = true
					allContent.WriteString("## SIZE LIMIT REACHED\nRemaining files were not added.\n")
					return nil
				}

				// Add file to total content
				allContent.WriteString(fileEntry)
				currentSize += fileEntrySize
				fileCount++
			} else if !config.SkipBinaryFiles {
				// For binary files, only add file information
				binaryInfo := fmt.Sprintf("## File: %s\n(Binary file, content skipped)\n\n", relPath)

				// Check size
				if currentSize+len(binaryInfo) <= config.MaxClipboardSize {
					allContent.WriteString(binaryInfo)
					currentSize += len(binaryInfo)
				} else {
					maxSizeReached = true
				}
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

// Check if directory is hidden
func shouldSkipDir(info os.FileInfo) bool {
	if !config.SkipHiddenDirs {
		return false
	}
	isHidden := info.IsDir() && len(info.Name()) > 0 && info.Name()[0] == '.'
	return isHidden
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

func readFileContent(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("Failed to read file '%s': %v", path, err)
	}
	return string(content), nil
}
