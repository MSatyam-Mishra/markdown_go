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
	client := flag.String("client", "", "The AI client to configure (claude, cursor)")
	flag.Parse()

	if *client == "" {
		fmt.Println("Error: --client flag is required (e.g., --client claude)")
		os.Exit(1)
	}

	configPath := getConfigPath(*client)
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
	client = strings.ToLower(client)

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
	case "cursor":
		return filepath.Join(home, ".cursor", "mcp.json")
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
	mcpServers, ok := config["mcpServers"].(map[string]interface{})
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
