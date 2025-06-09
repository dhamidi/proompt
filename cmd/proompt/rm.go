package main

import (
	"fmt"

	"github.com/dhamidi/proompt/pkg/picker"
	"github.com/dhamidi/proompt/pkg/prompt"
	"github.com/spf13/cobra"
)

// rmCmd creates the remove command
func rmCmd(manager prompt.Manager, pick picker.Picker) *cobra.Command {
	return &cobra.Command{
		Use:   "rm [name]",
		Short: "Remove a prompt",
		Long:  "Remove a prompt. If no name is provided, a picker will be used to select one.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var promptName string
			var err error

			// If a name is provided as argument
			if len(args) > 0 {
				promptName = args[0]
			} else {
				// No name provided, use picker to select existing prompt
				items, err := manager.GetAllForPicker()
				if err != nil {
					return fmt.Errorf("failed to get prompts for picker: %w", err)
				}

				if len(items) == 0 {
					return fmt.Errorf("no prompts found to remove")
				}

				selected, err := pick.Pick(items)
				if err != nil {
					return fmt.Errorf("picker failed: %w", err)
				}

				promptName = selected.Name
			}

			// Verify the prompt exists before attempting to delete
			promptInfo, err := manager.Get(promptName)
			if err != nil {
				if err == prompt.ErrPromptNotFound {
					return fmt.Errorf("prompt '%s' not found", promptName)
				}
				return fmt.Errorf("failed to get prompt: %w", err)
			}

			// Delete the prompt
			err = manager.Delete(promptName)
			if err != nil {
				return fmt.Errorf("failed to remove prompt: %w", err)
			}

			fmt.Printf("Removed prompt: %s (%s)\n", promptInfo.Name, promptInfo.Source)
			return nil
		},
	}
}
