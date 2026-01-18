package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fishfisher/apple-notes/internal/applescript"
	"github.com/fishfisher/apple-notes/internal/db"
	"github.com/spf13/cobra"
)

var (
	editTitle string
	editBody  string
	editForce bool
)

var editCmd = &cobra.Command{
	Use:   "edit [note-title]",
	Short: "Edit an existing note",
	Long: `Edit an existing note's title and/or body. Use --title to rename, --body to change content.

WARNING: Editing will replace the note body with plain text, which will remove:
  - Images and photos
  - Attachments and files
  - Tables and formatting
  - Sketches and drawings

Use --force to skip the confirmation prompt.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		noteTitle := args[0]

		// Verify note exists using SQLite
		database, err := db.Open()
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer database.Close()

		note, err := database.GetNote(noteTitle)
		if err != nil {
			return fmt.Errorf("note not found: %w", err)
		}

		// Use existing title if new title not provided
		newTitle := editTitle
		if newTitle == "" {
			newTitle = note.Title
		}

		// If body not provided via flag, read from stdin or prompt
		newBody := editBody
		if newBody == "" {
			fmt.Printf("Editing note '%s'\n", note.Title)
			fmt.Println("Enter new body (Ctrl+D when done):")
			scanner := bufio.NewScanner(os.Stdin)
			var lines []string
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}
			newBody = strings.Join(lines, "\n")
		}

		// Warn about attachments unless --force is used
		if !editForce {
			fmt.Println("\nWARNING: This will replace the note body with plain text.")
			fmt.Println("Any images, attachments, tables, or formatting will be lost.")
			fmt.Print("Continue? (y/N): ")
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Edit cancelled")
				return nil
			}
		}

		fmt.Printf("Updating note '%s'...\n", note.Title)
		if err := applescript.EditNote(note.Title, newTitle, newBody); err != nil {
			return fmt.Errorf("failed to edit note: %w", err)
		}

		fmt.Println("Note updated successfully")
		return nil
	},
}

func init() {
	editCmd.Flags().StringVarP(&editTitle, "title", "t", "", "New title for the note (optional)")
	editCmd.Flags().StringVarP(&editBody, "body", "b", "", "New body for the note (if not provided, reads from stdin)")
	editCmd.Flags().BoolVar(&editForce, "force", false, "Skip confirmation prompt (use with caution)")
}
