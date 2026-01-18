package cmd

import (
	"fmt"

	"github.com/fishfisher/apple-notes/internal/db"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show [note-id-or-title]",
	Short: "Show a specific note",
	Long:  `Display the full content of a specific note by ID or title.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		identifier := args[0]

		database, err := db.Open()
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer database.Close()

		note, err := database.GetNote(identifier)
		if err != nil {
			return err
		}

		fmt.Printf("Title:    %s\n", note.Title)
		fmt.Printf("Folder:   %s\n", note.Folder)
		fmt.Printf("Created:  %s\n", note.Created.Format("2006-01-02 15:04:05"))
		fmt.Printf("Modified: %s\n", note.Modified.Format("2006-01-02 15:04:05"))
		fmt.Printf("\n%s\n", note.Snippet)

		return nil
	},
}
