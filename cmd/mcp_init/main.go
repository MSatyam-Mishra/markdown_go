package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	client := flag.String("client", "", "The AI client to configure (claude, cursor, kiro, antigravity)")
	flag.Parse()

	if *client == "" {
		fmt.Println("Error: --client flag is required (e.g., --client claude)")
		os.Exit(1)
	}

	clientLower := strings.ToLower(*client)

	if clientLower == "antigravity" {
		err := injectAntigravitySkill()
		if err != nil {
			fmt.Printf("Failed to inject Antigravity skill: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Successfully installed MarkdownGo MCP server for Antigravity!")
		return
	}

	configPath := getConfigPath(clientLower)
	if configPath == "" {
		fmt.Printf("Error: Unsupported client '%s' or operating system.\n", *client)
		os.Exit(1)
	}

	fmt.Printf("Found config path: %s\n", configPath)

	err := injectMCPConfig(configPath)
	if err != nil {
		fmt.Printf("Failed to inject configuration: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully installed MarkdownGo MCP server for %s!\n", *client)
	fmt.Println("Please restart your AI client to apply the changes.")
}

func getConfigPath(client string) string {
	home, _ := os.UserHomeDir()

	switch client {
	case "claude":
		if runtime.GOOS == "windows" {
			appData := os.Getenv("APPDATA")
			return filepath.Join(appData, "Claude", "claude_desktop_config.json")
		} else if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "Application Support", "Claude", "claude_desktop_config.json")
		} else {
			// Linux/Other
			return filepath.Join(home, ".config", "Claude", "claude_desktop_config.json")
		}
	case "cursor", "windsurf":
		return filepath.Join(home, ".cursor", "mcp.json")
	case "kiro":
		return filepath.Join(home, ".kiro", "settings", "mcp.json")
	}
	return ""
}

func injectMCPConfig(configPath string) error {
	// Create directories if they don't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Read existing config
	var config map[string]interface{}
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			config = make(map[string]interface{})
		} else {
			return err
		}
	} else {
		if err := json.Unmarshal(data, &config); err != nil {
			return fmt.Errorf("failed to parse existing JSON: %v", err)
		}
	}

	// Ensure mcpServers object exists
	mcpServersAny, ok := config["mcpServers"]
	var mcpServers map[string]interface{}
	
	if ok {
		mcpServers, ok = mcpServersAny.(map[string]interface{})
	}
	if !ok {
		mcpServers = make(map[string]interface{})
		config["mcpServers"] = mcpServers
	}

	// Inject markdown_go
	mcpServers["markdown_go"] = map[string]interface{}{
		"command": "go",
		"args": []string{
			"run",
			"github.com/MSatyam-Mishra/markdown_go/cmd/mcp_server@latest",
		},
	}

	// Write back
	newData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, newData, 0644)
}

func injectAntigravitySkill() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	skillDir := filepath.Join(home, ".gemini", "config", "skills", "markdown_go_mcp")
	err = os.MkdirAll(skillDir, 0755)
	if err != nil {
		return err
	}

	skillPath := filepath.Join(skillDir, "SKILL.md")
	content := `---
name: markdown_go_mcp
description: "A native tool to convert local files (PDF, PPTX, DOCX, ZIP) and URLs (Webpages, Youtube Videos) into perfectly formatted markdown."
---

# markdown_go_mcp

This skill runs an MCP server that grants agents the ability to extract text and data from local files or websites and convert them perfectly into markdown.

## Usage

When an agent needs to read a PDF, Word document, Excel file, ZIP file, or fetch a Youtube Transcript/Web Page, they should call this MCP server.

**Server Name:** ` + "`markdown_go_mcp`" + `

**Tools Provided:**
- ` + "`convert_to_markdown(target: string)`" + `: Converts the target into markdown text.

## Configuration

To run this MCP server in Cursor, Claude Desktop, or Windsurf, configure your mcp.json to execute the built binary or run ` + "`go run ./cmd/mcp_server/main.go`" + `.
`

	err = os.WriteFile(skillPath, []byte(content), 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Created Antigravity skill at %s\n", skillPath)
	return nil
}
