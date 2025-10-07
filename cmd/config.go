package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management",
	Long:  "Manage mdcli configuration settings and preferences.",
}

var initConfigCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file",
	Long:  "Create a default configuration file with all available settings.",
	Run:   runInitConfig,
}

var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  "Display the current configuration settings from all sources.",
	Run:   runShowConfig,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(initConfigCmd)
	configCmd.AddCommand(showConfigCmd)
}

const defaultConfig = `# mdcli Configuration File
# This file contains default settings for mdcli

# Default theme for syntax highlighting
theme: dracula

# Default output format (terminal, html, pdf, text)
output_format: terminal

# Default terminal width for formatting
width: 80

# Enable automatic link detection
autolink: true

# Batch processing settings
batch:
  concurrent_workers: 4
  output_extension: .html
  recursive: true

# Server settings for live preview
serve:
  port: 8080
  bind: localhost
  auto_reload: true

# Rendering preferences
render:
  # Show progress bar for multiple files
  show_progress: true
  # Include metadata in output
  include_metadata: false

# Watch mode settings
watch:
  # Debounce delay in milliseconds
  debounce_delay: 100
  # Clear screen on update
  clear_screen: true

# Theme customizations (advanced users)
themes:
  # You can override specific theme colors here
  # custom:
  #   primary: "#ff0000"
  #   secondary: "#00ff00"

# File patterns to ignore in batch mode
ignore_patterns:
  - "node_modules/**"
  - ".git/**"
  - "*.tmp.md"`

func runInitConfig(cmd *cobra.Command, args []string) {
	// Determine config file path
	configPath := cfgFile
	if configPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding home directory: %v\n", err)
			os.Exit(1)
		}
		configPath = filepath.Join(home, ".mdcli.yaml")
	}

	// Check if config file already exists
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("Configuration file already exists at: %s\n", configPath)
		fmt.Print("Overwrite? (y/N): ")
		
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Configuration initialization cancelled.")
			return
		}
	}

	// Write default configuration
	err := os.WriteFile(configPath, []byte(defaultConfig), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing config file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Configuration file created at: %s\n", configPath)
	fmt.Println("You can now customize your settings by editing this file.")
	fmt.Println("💡 Use 'mdcli config show' to view current settings.")
}

func runShowConfig(cmd *cobra.Command, args []string) {
	fmt.Println("📋 Current Configuration:")
	fmt.Println("========================")

	if viper.ConfigFileUsed() != "" {
		fmt.Printf("Config file: %s\n\n", viper.ConfigFileUsed())
	} else {
		fmt.Println("Config file: <using defaults>")
	}

	// Display all current settings
	settings := map[string]interface{}{
		"Theme":         viper.GetString("theme"),
		"Output Format": viper.GetString("output_format"),
		"Width":         viper.GetInt("width"),
		"Autolink":      viper.GetBool("autolink"),
	}

	for key, value := range settings {
		fmt.Printf("%-15s: %v\n", key, value)
	}

	// Batch settings
	fmt.Println("\n🔄 Batch Processing:")
	fmt.Printf("%-15s: %d\n", "Workers", viper.GetInt("batch.concurrent_workers"))
	fmt.Printf("%-15s: %s\n", "Extension", viper.GetString("batch.output_extension"))
	fmt.Printf("%-15s: %t\n", "Recursive", viper.GetBool("batch.recursive"))

	// Server settings
	fmt.Println("\n🌐 Server:")
	fmt.Printf("%-15s: %d\n", "Port", viper.GetInt("serve.port"))
	fmt.Printf("%-15s: %s\n", "Bind", viper.GetString("serve.bind"))
	fmt.Printf("%-15s: %t\n", "Auto-reload", viper.GetBool("serve.auto_reload"))

	// Environment variables
	fmt.Println("\n🌍 Environment Variables:")
	envVars := []string{"MDCLI_THEME", "MDCLI_WIDTH", "MDCLI_FORMAT"}
	for _, env := range envVars {
		if value := os.Getenv(env); value != "" {
			fmt.Printf("%-15s: %s\n", env, value)
		}
	}

	fmt.Println("\n💡 You can override these settings using command-line flags.")
	fmt.Println("💡 Use 'mdcli config init' to create a configuration file.")
}