package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "proompt",
		Short: "A CLI tool for managing and using prompts",
		Long:  "Proompt is a CLI tool that helps you manage and use prompts with placeholder substitution.",
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
