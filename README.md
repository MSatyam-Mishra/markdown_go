# MarkdownGo

**MarkdownGo** is a lightweight, high-performance Go utility for converting various files to Markdown for use with LLMs and related text analysis pipelines. 

Inspired by Microsoft's Python-based *MarkItDown*, this library is built natively in Go. By leveraging Go's concurrency model (goroutines), MarkdownGo provides blazingly fast conversions—especially when extracting and converting files recursively from within ZIP archives.

MarkdownGo currently supports conversion from:

- **PDF** (Pure Go, no external C++ binaries required!)
- **PowerPoint** (.pptx)
- **Word** (.docx)
- **Excel**
- **Images** (EXIF metadata & OCR support)
- **HTML** (Local files)
- **Web URLs** (Scrapes generic HTML web pages perfectly into markdown)
- **Text-based formats** (CSV, JSON, XML)
- **ZIP files** (Iterates over all internal contents concurrently)
- **Youtube URLs** (Extracts video transcripts, metadata, and descriptions)
- **EPubs**


## Why Markdown?
Markdown is extremely close to plain text with minimal markup or formatting, yet it still provides a way to represent important document structure. Mainstream LLMs natively "speak" Markdown and often incorporate it into their responses unprompted. This suggests they have been trained on vast amounts of Markdown-formatted text and understand it well. As a side benefit, Markdown conventions are also highly token-efficient.

## Project Structure
MarkdownGo follows standard Go project layout conventions to ensure the codebase is highly scalable and easy to navigate for contributors:
- `markitdown.go`: The root package contains only the main orchestrator (`NewMarkItDown()`) to keep the public API incredibly clean.
- `pkg/converter/`: Contains all of the specialized data extraction logic and implementations (HTML, PDF, Youtube, etc.).
- `cmd/markdown_go/`: The entry point for building the global CLI tool.
- `example/`: A full-featured web app showcasing how to integrate MarkdownGo.

## Installation

### 1. As a Go Library
To use MarkdownGo inside your own Go applications:
```bash
go get github.com/MSatyam-Mishra/markdown_go
```

### 2. AI Tool Integration (MCP Server)

MarkdownGo includes a native **Model Context Protocol (MCP)** server, which allows AI agents (like Claude Desktop, Cursor, and Antigravity) to natively use this library as a tool for reading local files and scraping web pages!

#### Quick Start

Run the following command anywhere on your machine to automatically install and configure the MarkdownGo MCP server for your favorite AI client! 

*(No manual building or downloading required—Go will fetch it, compile it, and inject it into your client's configuration instantly!)*

```bash
# To install for Claude Desktop:
go run github.com/MSatyam-Mishra/markdown_go/cmd/mcp_init@latest --client claude

# To install for Cursor / Windsurf:
go run github.com/MSatyam-Mishra/markdown_go/cmd/mcp_init@latest --client cursor

# To install for Kiro IDE:
go run github.com/MSatyam-Mishra/markdown_go/cmd/mcp_init@latest --client kiro

# To install for Antigravity:
go run github.com/MSatyam-Mishra/markdown_go/cmd/mcp_init@latest --client antigravity
```

**Scope Options:**
By default, the installer configures the MCP server globally on your machine. You can pass the `--scope project` flag to install it only in your current workspace directory (e.g., creating a local `.agents` or `.cursor` configuration).

```bash
# Example: Install only for the current project
go run github.com/MSatyam-Mishra/markdown_go/cmd/mcp_init@latest --client cursor --scope project
```

Restart your AI client and ask it: *"Read the markdown from https://wikipedia.org"* and watch it use MarkdownGo seamlessly!

### 3. As a CLI Tool
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

Extract content from any generic website:
```bash
markdown_go -u "https://wikipedia.org/..."
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
	"github.com/MSatyam-Mishra/markdown_go/pkg/converter"
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

We also provide a beautiful, fully-functional Web UI frontend for testing conversions! 
The UI is built using modern **Shadcn UI** aesthetics, featuring a high-contrast design, clean rounded cards, and full **Dark Mode / Light Mode** support!

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

## API Documentation ¶

### Types ¶

#### type MarkItDown
```go
type MarkItDown struct
```
MarkItDown is the main orchestrator for conversions.

#### func NewMarkItDown
```go
func NewMarkItDown() *MarkItDown
```
NewMarkItDown initializes a new MarkItDown instance with all default converters registered (HTML, PDF, Word, Excel, ZIP, EPub, etc.).

#### func (*MarkItDown) ConvertFile
```go
func (m *MarkItDown) ConvertFile(ctx context.Context, path string) (string, error)
```
ConvertFile reads a local file and dynamically converts it to markdown based on its file extension.

#### func (*MarkItDown) ConvertURL
```go
func (m *MarkItDown) ConvertURL(ctx context.Context, url string) (string, error)
```
ConvertURL scrapes a webpage or YouTube video and extracts its contents and transcript perfectly into markdown.

---

### Package `converter` ¶

#### type Converter
```go
type Converter interface {
	Convert(ctx context.Context, r io.Reader, opts *Options) (string, error)
}
```
Converter is the interface that wraps the Convert method for specific format implementations.

#### type ConverterAPI
```go
type ConverterAPI interface {
	Convert(ctx context.Context, r io.Reader, opts *Options) (string, error)
}
```
ConverterAPI is the core orchestrator interface that allows plugins (like ZipConverter) to recursively convert embedded files without causing import cycles.

#### type Options
```go
type Options struct {
	Extension string
	FileName  string
	URL       string
}
```
Options provides metadata and configuration for the conversion process.
