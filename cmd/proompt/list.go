package main

import (
	"fmt"

	"github.com/dhamidi/proompt/pkg/prompt"
	"github.com/spf13/cobra"
)

// listCmd creates the list command
func listCmd(manager prompt.Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all available prompts",
		Long:  "List all available prompts from all configured locations (directory, project, project-local, user)",
		RunE: func(cmd *cobra.Command, args []string) error {
			prompts, err := manager.List()
			if err != nil {
				return fmt.Errorf("failed to list prompts: %w", err)
			}

			if len(prompts) == 0 {
				fmt.Println("No prompts found")
				return nil
			}

			fmt.Printf("Found %d prompt(s):\n\n", len(prompts))

			for _, prompt := range prompts {
				fmt.Printf("%-20s %-15s %s\n", prompt.Name, fmt.Sprintf("(%s)", prompt.Source), prompt.Path)
			}

			return nil
		},
	}
}
