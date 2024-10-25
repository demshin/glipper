# Glipper

Glipper is a small utility tool written in Go that recursively collects the content of all files in a specified directory and copies it to the clipboard. It is intended for developers and anyone who needs a quick way to gather and share file contents.

## Features

Recursively traverses directories and reads all file contents.

Skips hidden directories (e.g., .git or .idea).

Copies all collected content to the clipboard.

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

## Requirements

Go 1.16 or higher

golang.design/x/clipboard library for clipboard support

## Notes

Glipper skips hidden directories (those starting with a dot .) to avoid collecting content from unnecessary system or configuration folders.

The content of each file is separated by a divider for easier readability.

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Contributing

Feel free to submit issues and pull requests if you find any bugs or have suggestions for improvements.

## Future Improvements

Add support for filtering files by extension.

Implement better error handling for large directories or large file sizes.

## Author

Created by Aleksandr Demshin.

Website: https://demsh.in

Email: aleksandr@demsh.in

Telegram: @demshin

