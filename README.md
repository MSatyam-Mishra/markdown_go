# MarkdownGo

**MarkdownGo** is a lightweight, high-performance Go utility for converting various files to Markdown for use with LLMs and related text analysis pipelines. 

Inspired by Microsoft's Python-based *MarkItDown*, this library is built natively in Go. By leveraging Go's concurrency model (goroutines), MarkdownGo provides blazingly fast conversions—especially when extracting and converting files recursively from within ZIP archives.

MarkdownGo currently supports conversion from:

- **PDF** (Pure Go, no external C++ binaries required!)
- **PowerPoint** (.pptx)
- **Word** (.docx)
- **Excel**
- **Images** (EXIF metadata & OCR support)
- **HTML**
- **Text-based formats** (CSV, JSON, XML)
- **ZIP files** (Iterates over all internal contents concurrently)
- **Youtube URLs** (Extracts video metadata and descriptions)
- **EPubs**
- ... and more!

## Why Markdown?
Markdown is extremely close to plain text, with minimal markup or formatting, but still provides a way to represent important document structure. Mainstream LLMs natively "speak" Markdown, and often incorporate Markdown into their responses unprompted. This suggests that they have been trained on vast amounts of Markdown-formatted text, and understand it well. As a side benefit, Markdown conventions are also highly token-efficient.

## Installation

### 1. As a Go Library
To use MarkdownGo inside your own Go applications:
```bash
go get github.com/MSatyam-Mishra/markdown_go
```

### 2. As a CLI Tool
To install the MarkdownGo command-line interface globally on your machine:
```bash
go install github.com/MSatyam-Mishra/markdown_go/cmd/markdown_go@latest
```

## Usage

### Command-Line (CLI)

Convert a local file and output to stdout:
```bash
markdown_go -i path-to-file.pdf
```

Save the conversion directly to a Markdown file:
```bash
markdown_go -i path-to-file.pdf -o document.md
```

Extract metadata and text from a YouTube URL:
```bash
markdown_go -u "https://youtube.com/watch?v=..."
```

### Go API

Basic usage in your Go code:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/MSatyam-Mishra/markdown_go"
)

func main() {
	// Initialize the converter
	md := markdown_go.NewMarkItDown()
	ctx := context.Background()
	
	// Convert a local file
	result, err := md.ConvertFile(ctx, "document.pdf")
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println(result)
}
```

Convert a YouTube URL directly:

```go
result, err := md.ConvertURL(ctx, "https://youtube.com/watch?v=...")
```

## Example Web App

We have also provided a beautiful, fully-functional Web UI frontend to test out conversions!

To run it locally:
```bash
cd example
go run main.go
```
Then navigate to **http://localhost:8080** in your browser.

## Security Considerations

MarkdownGo performs I/O with the privileges of the current process. Like `os.Open()`, it will access resources that the process itself can access.

**Sanitize your inputs:** Do not pass untrusted input directly to MarkdownGo. If any part of the input may be controlled by an untrusted user or system, such as in hosted or server-side applications, it must be validated and restricted before calling MarkdownGo.

## Contributing

This project welcomes contributions and suggestions. Feel free to submit PRs, add support for more APIs (like Azure Document Intelligence or OpenAI Whisper), or expand the URL conversion capabilities!
