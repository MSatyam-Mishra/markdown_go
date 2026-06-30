package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/MSatyam-Mishra/markdown_go"
)

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"MarkdownGo Server",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// Define the tool
	convertTool := mcp.NewTool("convert_to_markdown",
		mcp.WithDescription("Converts a local file (PDF, PPTX, DOCX, Images, ZIP) or a URL (Web Page, YouTube video) to markdown text."),
		mcp.WithString("target",
			mcp.Required(),
			mcp.Description("The absolute file path or URL to convert to markdown."),
		),
	)

	// Add the tool to the server
	s.AddTool(convertTool, convertHandler)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("MCP Server error: %v\n", err)
	}
}

// Handler function for the convert_to_markdown tool
func convertHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	targetAny, exists := request.Params.Arguments.(map[string]interface{})["target"]
	if !exists {
		return mcp.NewToolResultError("target argument is required"), nil
	}
	
	target, ok := targetAny.(string)
	if !ok || target == "" {
		return mcp.NewToolResultError("target is required and must be a string"), nil
	}

	// Initialize MarkdownGo
	md := markdown_go.NewMarkItDown()
	
	var result string
	var convertErr error

	if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
		result, convertErr = md.ConvertURL(ctx, target)
	} else {
		result, convertErr = md.ConvertFile(ctx, target)
	}

	if convertErr != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to convert %s: %v", target, convertErr)), nil
	}

	return mcp.NewToolResultText(result), nil
}
