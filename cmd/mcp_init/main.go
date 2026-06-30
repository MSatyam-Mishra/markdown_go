package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

func main() {
	client := flag.String("client", "", "The AI client to configure (claude, cursor, kiro, antigravity)")
	scope := flag.String("scope", "", "The installation scope (global or project)")
	flag.Parse()

	if *client == "" {
		fmt.Println("Error: --client flag is required (e.g., --client claude)")
		os.Exit(1)
	}

	if *scope == "" {
		prompt := &survey.Select{
			Message: "Select installation scope:",
			Options: []string{"Global (Recommended)", "Project (Current Directory Only)"},
		}
		var choice string
		err := survey.AskOne(prompt, &choice)
		if err != nil {
			fmt.Printf("Prompt failed: %v\n", err)
			os.Exit(1)
		}
		if strings.HasPrefix(choice, "Project") {
			*scope = "project"
		} else {
			*scope = "global"
		}
	}

	clientLower := strings.ToLower(*client)
	scopeLower := strings.ToLower(*scope)

	if scopeLower != "global" && scopeLower != "project" {
		fmt.Println("Error: --scope flag must be either 'global' or 'project'")
		os.Exit(1)
	}

	if clientLower == "antigravity" {
		err := injectAntigravitySkill(scopeLower)
		if err != nil {
			fmt.Printf("Failed to inject Antigravity skill: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully installed MarkdownGo MCP server for Antigravity (%s scope)!\n", scopeLower)
		return
	}

	configPath := getConfigPath(clientLower, scopeLower)
	if configPath == "" {
		fmt.Printf("Error: Unsupported client '%s' or scope '%s'.\n", *client, scopeLower)
		os.Exit(1)
	}

	fmt.Printf("Found config path: %s\n", configPath)

	err := injectMCPConfig(configPath)
	if err != nil {
		fmt.Printf("Failed to inject configuration: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully installed MarkdownGo MCP server for %s (%s scope)!\n", *client, scopeLower)
	fmt.Println("Please restart your AI client to apply the changes.")
}

func getConfigPath(client, scope string) string {
	home, _ := os.UserHomeDir()
	pwd, _ := os.Getwd()

	switch client {
	case "claude":
		if scope == "project" {
			fmt.Println("Warning: Claude Desktop does not support project-level configurations natively. Installing globally instead.")
		}
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
		if scope == "project" {
			return filepath.Join(pwd, ".cursor", "mcp.json")
		}
		return filepath.Join(home, ".cursor", "mcp.json")
	case "kiro":
		if scope == "project" {
			return filepath.Join(pwd, ".kiro", "settings", "mcp.json")
		}
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

func injectAntigravitySkill(scope string) error {
	var skillDir string
	if scope == "global" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		skillDir = filepath.Join(home, ".gemini", "config", "skills", "markdown_go_mcp")
	} else {
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		skillDir = filepath.Join(pwd, ".agents", "skills", "markdown_go_mcp")
	}

	err := os.MkdirAll(skillDir, 0755)
	if err != nil {
		return err
	}

	skillPath := filepath.Join(skillDir, "SKILL.md")
	content := `---
name: markdown_go_mcp
description: "Native tool to convert local files (PDF, PPTX, DOCX, ZIP) and URLs (Webpages, Youtube Videos) into markdown."
---

# markdown_go_mcp

This skill provides a Model Context Protocol (MCP) server that extracts text and data from local files or websites and converts them into markdown.

## When to use this skill
Use this skill whenever the user asks you to:
- Read or extract text from a **PDF, Word document (.docx), PowerPoint (.pptx), or Excel file**.
- Scrape or read a **web page** (URL).
- Extract a transcript and metadata from a **YouTube video** (URL).
- Read contents recursively from a **ZIP file**.

## How to use this tool

The MCP server exposes a tool called ` + "`convert_to_markdown`" + `.

### Arguments
The tool requires a single argument:
- ` + "`target`" + ` (string): This must be either a **fully qualified URL** (e.g. "https://en.wikipedia.org/...") or an **absolute local file path** (e.g. "C:/Users/name/document.pdf").

### Execution
Depending on how the server is registered in your environment, you can call it directly if it's eager-loaded (` + "`markdown_go_convert_to_markdown`" + `), or use the ` + "`call_mcp_tool`" + ` tool with:
- **ServerName**: ` + "`markdown_go`" + `
- **ToolName**: ` + "`convert_to_markdown`" + `
- **Arguments**: ` + "`{\"target\": \"<url_or_absolute_path>\"}`" + `

### Handling the Output
The tool will return the raw markdown string. 
1. **Do not** attempt to summarize the entire text if the user wants an exact translation.
2. If the text is extremely large, use it to answer the user's specific questions rather than printing it all out.
`

	err = os.WriteFile(skillPath, []byte(content), 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Created detailed Antigravity skill at %s\n", skillPath)
	return nil
}
