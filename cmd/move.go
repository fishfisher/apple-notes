package cmd

import (
	"fmt"

	"github.com/fishfisher/apple-notes/internal/applescript"
	"github.com/fishfisher/apple-notes/internal/db"
	"github.com/spf13/cobra"
)

var moveCmd = &cobra.Command{
	Use:   "move [note-title] [target-folder]",
	Short: "Move a note to a different folder",
	Long:  `Move an existing note to a different folder.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		noteTitle := args[0]
		targetFolder := args[1]

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

		if note.Folder == targetFolder {
			fmt.Printf("Note '%s' is already in folder '%s'\n", note.Title, targetFolder)
			return nil
		}

		fmt.Printf("Moving note '%s' from '%s' to '%s'...\n", note.Title, note.Folder, targetFolder)
		if err := applescript.MoveNote(note.Title, targetFolder); err != nil {
			return fmt.Errorf("failed to move note: %w", err)
		}

		fmt.Println("Note moved successfully")
		return nil
	},
}
