# Glipper

Glipper is a small utility tool written in Go that recursively collects the content of all files in a specified directory and copies it to the clipboard. It is intended for developers and anyone who needs a quick way to gather and share file contents.

## Requirements

- macOS (supports both Intel and ARM architectures)
- Go 1.24.1 or higher (for building from source)
- github.com/atotto/clipboard library for clipboard support

## Features

- Recursively traverses directories and reads all file contents.
- Skips hidden directories (e.g., .git or .idea).
- Copies all collected content to the clipboard.
- Automatically limits clipboard size to prevent issues with large directories.
- Configurable via simple config file and command-line options.

## Installation

To install Glipper, you need to have Go installed on your system.

Clone the repository:

```bash
git clone https://github.com/yourusername/glipper.git
cd glipper
```

Build the project:

```bash
go build -o glipper main.go
```

## Usage

To use Glipper, run the compiled executable with the path to the directory you want to collect content from:

`./glipper /path/to/your/dir`

For example:

`./glipper ~/myproject`

This will recursively collect the content of all files in the specified directory and copy it to your clipboard.

### Command-line Options

Glipper supports the following command-line options:

- `-size=64000`: Maximum clipboard size in bytes (default: 64KB)
- `-skip-binary=true`: Skip binary files (default: true)
- `-skip-hidden=true`: Skip hidden directories (default: true)
- `-help`: Show help information

Example with options:

```bash
./glipper -size=100000 -skip-binary=false ~/myproject
```

## Configuration

Glipper can be configured via a configuration file located at `~/.config/glipper/.glipper.conf`. The file uses a simple key=value format:

```
# Glipper configuration file
# Format: key=value

max_clipboard_size=64000
skip_binary_files=true
skip_hidden_dirs=true
```

The configuration file is automatically created with default values on first run. Any changes made via command-line options will be saved to this file for future use.

## Notes

- Glipper skips hidden directories (those starting with a dot .) to avoid collecting content from unnecessary system or configuration folders.
- The content of each file is separated by a divider for easier readability.
- Files larger than 1MB are automatically skipped to prevent memory issues.
- Content is automatically truncated if it exceeds the configured maximum clipboard size.

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Contributing

Feel free to submit issues and pull requests if you find any bugs or have suggestions for improvements.

## Future Improvements

- Add support for filtering files by extension.
- Implement better error handling for large directories or large file sizes.
- Add support for configurable file size limits.
- Add more output format options.

## Author

Created by Aleksandr Demshin.

Website: https://demsh.in

Email: aleksandr@demsh.in

Telegram: @demshin
