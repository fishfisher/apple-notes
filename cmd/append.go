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
	appendContent string
	appendByTitle bool
)

var appendCmd = &cobra.Command{
	Use:   "append [note-id-or-title]",
	Short: "Append content to an existing note",
	Long:  `Append content to an existing note without replacing the current content.

By default, numeric input is treated as a note ID. Use --by-title to search by title instead.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		identifier := args[0]

		// Verify note exists
		database, err := db.Open()
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer database.Close()

		var note *db.Note
		if appendByTitle {
			note, err = database.GetNoteByTitle(identifier)
		} else {
			note, err = database.GetNote(identifier)
		}
		if err != nil {
			return fmt.Errorf("note not found: %w", err)
		}

		// If content not provided via flag, read from stdin
		content := appendContent
		if content == "" {
			fmt.Println("Enter content to append (Ctrl+D when done):")
			scanner := bufio.NewScanner(os.Stdin)
			var lines []string
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}
			content = strings.Join(lines, "\n")
		}

		fmt.Printf("Appending to note '%s'...\n", note.Title)
		if err := applescript.AppendNote(note.Title, content); err != nil {
			return fmt.Errorf("failed to append to note: %w", err)
		}

		fmt.Println("Content appended successfully")
		return nil
	},
}

func init() {
	appendCmd.Flags().StringVarP(&appendContent, "content", "c", "", "Content to append (if not provided, reads from stdin)")
	appendCmd.Flags().BoolVarP(&appendByTitle, "by-title", "t", false, "Search by title (even if input is numeric)")
}
