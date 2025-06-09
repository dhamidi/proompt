package main

import (
	"fmt"

	"github.com/dhamidi/proompt/pkg/prompt"
	"github.com/spf13/cobra"
)

// showCmd creates the show command
func showCmd(manager prompt.Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "show <name>",
		Short: "Show a specific prompt",
		Long:  "Show the content of a specific prompt by name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			
			promptInfo, err := manager.Get(name)
			if err != nil {
				return fmt.Errorf("failed to get prompt '%s': %w", name, err)
			}
			
			if promptInfo == nil {
				return fmt.Errorf("prompt '%s' not found", name)
			}

			fmt.Printf("Name: %s\n", promptInfo.Name)
			fmt.Printf("Source: %s\n", promptInfo.Source)
			fmt.Printf("Path: %s\n", promptInfo.Path)
			fmt.Printf("\nContent:\n%s\n", promptInfo.Content)

			return nil
		},
	}
}
