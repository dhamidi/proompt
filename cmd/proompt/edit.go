package main

import (
	"fmt"

	"github.com/dhamidi/proompt/pkg/editor"
	"github.com/dhamidi/proompt/pkg/picker"
	"github.com/dhamidi/proompt/pkg/prompt"
	"github.com/spf13/cobra"
)

// editCmd creates the edit command
func editCmd(manager prompt.Manager, pick picker.Picker, ed editor.Editor) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit [name]",
		Short: "Edit a prompt",
		Long:  "Edit a prompt. If no name is provided, a picker will be used to select one. Use location flags to create new prompts at specific levels.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get location flags
			directory, _ := cmd.Flags().GetBool("directory")
			project, _ := cmd.Flags().GetBool("project")
			projectLocal, _ := cmd.Flags().GetBool("project-local")
			user, _ := cmd.Flags().GetBool("user")

			// Count how many location flags are set
			locationCount := 0
			var targetLocation string
			if directory {
				locationCount++
				targetLocation = "directory"
			}
			if project {
				locationCount++
				targetLocation = "project"
			}
			if projectLocal {
				locationCount++
				targetLocation = "project-local"
			}
			if user {
				locationCount++
				targetLocation = "user"
			}

			// Validate that only one location flag is set
			if locationCount > 1 {
				return fmt.Errorf("only one location flag can be specified")
			}

			var promptName string
			var promptInfo *prompt.PromptInfo
			var err error

			// If a name is provided as argument
			if len(args) > 0 {
				promptName = args[0]
				
				// Try to get existing prompt
				promptInfo, err = manager.Get(promptName)
				if err != nil && err != prompt.ErrPromptNotFound {
					return fmt.Errorf("failed to get prompt: %w", err)
				}

				// If prompt doesn't exist and no location flag specified, error
				if err == prompt.ErrPromptNotFound && locationCount == 0 {
					return fmt.Errorf("prompt '%s' not found. Use a location flag (--directory, --project, --project-local, --user) to create it", promptName)
				}

				// If prompt exists and location flag specified, error
				if err == nil && locationCount > 0 {
					return fmt.Errorf("prompt '%s' already exists at %s. Cannot specify location flag for existing prompts", promptName, promptInfo.Source)
				}
			} else {
				// No name provided, use picker to select existing prompt
				if locationCount > 0 {
					return fmt.Errorf("cannot use location flags without providing a prompt name")
				}

				items, err := manager.GetAllForPicker()
				if err != nil {
					return fmt.Errorf("failed to get prompts for picker: %w", err)
				}

				if len(items) == 0 {
					return fmt.Errorf("no prompts found to edit")
				}

				selected, err := pick.Pick(items)
				if err != nil {
					return fmt.Errorf("picker failed: %w", err)
				}

				promptName = selected.Name
				promptInfo, err = manager.Get(promptName)
				if err != nil {
					return fmt.Errorf("failed to get selected prompt: %w", err)
				}
			}

			// Handle creating new prompt
			if promptInfo == nil {
				// Create new prompt with empty content
				err = manager.Create(promptName, "", targetLocation)
				if err != nil {
					return fmt.Errorf("failed to create prompt: %w", err)
				}

				// Get the newly created prompt
				promptInfo, err = manager.Get(promptName)
				if err != nil {
					return fmt.Errorf("failed to get newly created prompt: %w", err)
				}
			}

			// Edit the prompt file
			err = ed.Edit(promptInfo.Path)
			if err != nil {
				return fmt.Errorf("editor failed: %w", err)
			}

			fmt.Printf("Edited prompt: %s (%s)\n", promptInfo.Name, promptInfo.Source)
			return nil
		},
	}

	// Add location flags
	cmd.Flags().Bool("directory", false, "Create prompt in directory level (./prompts/)")
	cmd.Flags().Bool("project", false, "Create prompt in project level (project root/prompts/)")
	cmd.Flags().Bool("project-local", false, "Create prompt in project-local level (.git/info/prompts/)")
	cmd.Flags().Bool("user", false, "Create prompt in user level (config directory)")

	return cmd
}
