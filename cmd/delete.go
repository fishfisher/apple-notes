package cmd

import (
	"fmt"

	"github.com/fishfisher/apple-notes/internal/applescript"
	"github.com/fishfisher/apple-notes/internal/db"
	"github.com/spf13/cobra"
)

var (
	deleteForce   bool
	deleteByTitle bool
)

var deleteCmd = &cobra.Command{
	Use:   "delete [note-id-or-title]",
	Short: "Delete a note",
	Long:  `Delete a note. Use --force to skip confirmation.

By default, numeric input is treated as a note ID. Use --by-title to search by title instead.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		identifier := args[0]

		// Verify note exists using SQLite
		database, err := db.Open()
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer database.Close()

		var note *db.Note
		if deleteByTitle {
			note, err = database.GetNoteByTitle(identifier)
		} else {
			note, err = database.GetNote(identifier)
		}
		if err != nil {
			return fmt.Errorf("note not found: %w", err)
		}

		// Confirm deletion unless --force is used
		if !deleteForce {
			fmt.Printf("Delete note '%s' from folder '%s'? (y/N): ", note.Title, note.Folder)
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Deletion cancelled")
				return nil
			}
		}

		fmt.Printf("Deleting note '%s'...\n", note.Title)
		if err := applescript.DeleteNote(note.Title); err != nil {
			return fmt.Errorf("failed to delete note: %w", err)
		}

		fmt.Println("Note deleted successfully")
		return nil
	},
}

func init() {
	deleteCmd.Flags().BoolVarP(&deleteForce, "force", "y", false, "Skip confirmation prompt")
	deleteCmd.Flags().BoolVarP(&deleteByTitle, "by-title", "t", false, "Search by title (even if input is numeric)")
}
