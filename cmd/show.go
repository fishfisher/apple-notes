package cmd

import (
	"fmt"

	"github.com/fishfisher/apple-notes/internal/applescript"
	"github.com/fishfisher/apple-notes/internal/db"
	"github.com/spf13/cobra"
)

var showByTitle bool

var showCmd = &cobra.Command{
	Use:   "show [note-id-or-title]",
	Short: "Show a specific note",
	Long:  `Display the full content of a specific note.

By default, numeric input is treated as a note ID. Use --by-title to search by title instead.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		identifier := args[0]

		database, err := db.Open()
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer database.Close()

		var note *db.Note
		if showByTitle {
			note, err = database.GetNoteByTitle(identifier)
		} else {
			note, err = database.GetNote(identifier)
		}
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

func init() {
	showCmd.Flags().BoolVarP(&showByTitle, "by-title", "t", false, "Search by title (even if input is numeric)")
}
