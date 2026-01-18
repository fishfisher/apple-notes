package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fishfisher/apple-notes/internal/applescript"
	"github.com/spf13/cobra"
)

var (
	addFolder string
	addBody   string
)

var addCmd = &cobra.Command{
	Use:   "add [title]",
	Short: "Add a new note",
	Long:  `Create a new note with the specified title and body in the given folder.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := args[0]

		// If body not provided via flag, read from stdin or prompt
		body := addBody
		if body == "" {
			fmt.Println("Enter note body (Ctrl+D when done):")
			scanner := bufio.NewScanner(os.Stdin)
			var lines []string
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}
			body = strings.Join(lines, "\n")
		}

		// Default to "Notes" folder if not specified
		folder := addFolder
		if folder == "" {
			folder = "Notes"
		}

		fmt.Printf("Creating note '%s' in folder '%s'...\n", title, folder)
		if err := applescript.AddNote(title, body, folder); err != nil {
			return fmt.Errorf("failed to add note: %w", err)
		}

		fmt.Println("Note created successfully")
		return nil
	},
}

func init() {
	addCmd.Flags().StringVarP(&addFolder, "folder", "f", "Notes", "Folder to create the note in")
	addCmd.Flags().StringVarP(&addBody, "body", "b", "", "Note body (if not provided, reads from stdin)")
}
