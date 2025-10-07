package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool
	version = "2.0.0"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mdcli",
	Short: "A powerful Markdown CLI processor and renderer",
	Long: `mdcli is a feature-rich command-line tool for processing Markdown files.
It supports multiple output formats, themes, live preview, batch processing,
and many other advanced features to enhance your Markdown workflow.

Usage:
  mdcli [file.md]                 # Render file directly (legacy mode)
  mdcli render [files...]         # Explicit render command
  mdcli serve [file.md]           # Start live preview server
  mdcli watch [file.md]           # Watch for changes
  mdcli batch [directory]         # Process multiple files`,
	Version: version,
	Args:    cobra.ArbitraryArgs,
	Run:     runRootCommand,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// runRootCommand handles direct file rendering for backward compatibility
func runRootCommand(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		// No arguments provided, show help
		cmd.Help()
		return
	}

	// Check if any arguments look like subcommands
	for _, arg := range args {
		if isSubcommand(arg) {
			// If the first argument is a known subcommand, let cobra handle it normally
			// This shouldn't happen due to cobra's parsing, but just in case
			fmt.Fprintf(os.Stderr, "Unknown command '%s'. Use 'mdcli help' for available commands.\n", arg)
			os.Exit(1)
		}
	}

	// All arguments should be files, delegate to render command
	if verbose {
		fmt.Fprintf(os.Stderr, "Legacy mode: rendering files directly\n")
	}

	// Get flag values from the root command and set them for render
	outputFile, _ := cmd.Flags().GetString("output")
	outputFormat, _ := cmd.Flags().GetString("format")
	theme, _ := cmd.Flags().GetString("theme")
	width, _ := cmd.Flags().GetInt("width")
	autolink, _ := cmd.Flags().GetBool("autolink")
	showProgress, _ := cmd.Flags().GetBool("progress")

	// Call render function with a simulated render command context
	runRenderWithFlags(cmd, args, outputFile, outputFormat, theme, width, autolink, showProgress)
}

// isSubcommand checks if a string matches any known subcommand
func isSubcommand(arg string) bool {
	subcommands := []string{"render", "serve", "watch", "batch", "interactive", "themes", "config", "help", "completion"}
	for _, sub := range subcommands {
		if arg == sub {
			return true
		}
	}
	return false
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mdcli.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Add common render flags to root command for backward compatibility
	rootCmd.Flags().StringP("output", "o", "", "Output file path")
	rootCmd.Flags().StringP("format", "f", "terminal", "Output format (terminal, html, pdf, text)")
	rootCmd.Flags().StringP("theme", "t", "", "Syntax highlighting theme")
	rootCmd.Flags().IntP("width", "w", 0, "Terminal width for formatting")
	rootCmd.Flags().Bool("autolink", true, "Enable automatic link detection")
	rootCmd.Flags().Bool("progress", false, "Show progress bar")

	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("output", rootCmd.Flags().Lookup("output"))
	viper.BindPFlag("format", rootCmd.Flags().Lookup("format"))
	viper.BindPFlag("theme", rootCmd.Flags().Lookup("theme"))
	viper.BindPFlag("width", rootCmd.Flags().Lookup("width"))
	viper.BindPFlag("autolink", rootCmd.Flags().Lookup("autolink"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".mdcli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".mdcli")
	}

	// Set environment variable prefix
	viper.SetEnvPrefix("MDCLI")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil && verbose {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	// Set default values
	viper.SetDefault("theme", "dracula")
	viper.SetDefault("width", 80)
	viper.SetDefault("output_format", "terminal")
	viper.SetDefault("autolink", true)
	viper.SetDefault("batch.concurrent_workers", 4)
	viper.SetDefault("batch.output_extension", ".html")
	viper.SetDefault("batch.recursive", true)
	viper.SetDefault("serve.port", 8080)
	viper.SetDefault("serve.bind", "localhost")
	viper.SetDefault("serve.auto_reload", true)
	viper.SetDefault("render.show_progress", true)
	viper.SetDefault("render.include_metadata", false)
	viper.SetDefault("watch.debounce_delay", 100)
	viper.SetDefault("watch.clear_screen", true)
}
