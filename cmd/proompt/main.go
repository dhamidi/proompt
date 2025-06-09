package main

import (
	"fmt"
	"os"

	"github.com/dhamidi/proompt/pkg/config"
	"github.com/dhamidi/proompt/pkg/copier"
	"github.com/dhamidi/proompt/pkg/editor"
	"github.com/dhamidi/proompt/pkg/filesystem"
	"github.com/dhamidi/proompt/pkg/picker"
	"github.com/dhamidi/proompt/pkg/prompt"
	"github.com/spf13/cobra"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Get current working directory for filesystem root
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get current directory: %v\n", err)
		os.Exit(1)
	}

	// Initialize components
	fs := filesystem.NewRealFilesystem(cwd)
	resolver := prompt.NewDefaultLocationResolver(fs)
	manager := prompt.NewDefaultManager(fs, resolver)
	pick := picker.NewRealPicker(cfg.Picker)
	ed := editor.NewRealEditor(cfg.Editor)
	parser := prompt.NewDefaultParser()
	
	// Get copy command from environment or use default
	copyCommand := os.Getenv("PROOMPT_COPY_COMMAND")
	if copyCommand == "" {
		copyCommand = "pbcopy"
	}
	cop := copier.NewRealCopier(copyCommand)

	// Create root command
	rootCmd := &cobra.Command{
		Use:   "proompt",
		Short: "A CLI tool for managing and using prompts",
		Long:  "Proompt is a CLI tool that helps you manage and use prompts with placeholder substitution.",
	}

	// Add subcommands
	rootCmd.AddCommand(
		listCmd(manager),
		showCmd(manager),
		editCmd(manager, pick, ed),
		rmCmd(manager, pick),
		pickCmd(manager, pick, ed, parser, fs, cop),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
