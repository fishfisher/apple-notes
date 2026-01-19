package cmd

import (
	"fmt"

	"github.com/fishfisher/apple-notes/internal/applescript"
	"github.com/fishfisher/apple-notes/internal/db"
	"github.com/spf13/cobra"
)

var moveByTitle bool

var moveCmd = &cobra.Command{
	Use:   "move [note-id-or-title] [target-folder]",
	Short: "Move a note to a different folder",
	Long:  `Move an existing note to a different folder.

By default, numeric input is treated as a note ID. Use --by-title to search by title instead.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		identifier := args[0]
		targetFolder := args[1]

		// Verify note exists using SQLite
		database, err := db.Open()
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer database.Close()

		var note *db.Note
		if moveByTitle {
			note, err = database.GetNoteByTitle(identifier)
		} else {
			note, err = database.GetNote(identifier)
		}
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

func init() {
	moveCmd.Flags().BoolVarP(&moveByTitle, "by-title", "t", false, "Search by title (even if input is numeric)")
}
