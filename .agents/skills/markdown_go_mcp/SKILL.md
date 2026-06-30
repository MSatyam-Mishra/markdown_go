---
name: markdown_go_mcp
description: "A native tool to convert local files (PDF, PPTX, DOCX, ZIP) and URLs (Webpages, Youtube Videos) into perfectly formatted markdown."
---

# markdown_go_mcp

This skill runs an MCP server that grants agents the ability to extract text and data from local files or websites and convert them perfectly into markdown.

## Usage

When an agent needs to read a PDF, Word document, Excel file, ZIP file, or fetch a Youtube Transcript/Web Page, they should call this MCP server.

**Server Name:** `markdown_go_mcp`

**Tools Provided:**
- `convert_to_markdown(target: string)`: Converts the `target` (which must be an absolute file path or a URL) into markdown text.

## Configuration

To run this MCP server in Cursor, Claude Desktop, or Windsurf, configure your `mcp.json` to execute the built binary or run `go run ./cmd/mcp_server/main.go`.
