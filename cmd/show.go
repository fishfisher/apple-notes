package cmd

import (
	"fmt"

	"github.com/fishfisher/apple-notes/internal/applescript"
	"github.com/fishfisher/apple-notes/internal/db"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show [note-id]",
	Short: "Show a specific note",
	Long:  `Display the full content of a specific note by ID.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		noteID := args[0]

		database, err := db.Open()
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer database.Close()

		note, err := database.GetNote(noteID)
		if err != nil {
			return err
		}

		// Check for duplicate titles
		count, err := database.CountNotesByTitle(note.Title)
		if err == nil && count > 1 {
			fmt.Printf("⚠️  Warning: Found %d notes with title '%s'. Showing most recently modified.\n\n", count, note.Title)
		}

		fmt.Printf("ID:       %s\n", note.ID)
		fmt.Printf("Title:    %s\n", note.Title)
		fmt.Printf("Folder:   %s\n", note.Folder)
		fmt.Printf("Created:  %s\n", note.Created.Format("2006-01-02 15:04:05"))
		fmt.Printf("Modified: %s\n", note.Modified.Format("2006-01-02 15:04:05"))
		fmt.Printf("\n")

		// Get full body from AppleScript instead of just snippet
		body, err := applescript.GetNoteBody(note.Title)
		if err != nil {
			// Fallback to snippet if AppleScript fails
			fmt.Printf("Warning: Could not retrieve full note body, showing snippet only: %v\n\n", err)
			fmt.Printf("%s\n", note.Snippet)
		} else {
			fmt.Printf("%s\n", body)
		}

		return nil
	},
}
