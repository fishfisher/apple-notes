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

var appendContent string

var appendCmd = &cobra.Command{
	Use:   "append [note-title]",
	Short: "Append content to an existing note",
	Long:  `Append content to an existing note without replacing the current content.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		noteTitle := args[0]

		// Verify note exists
		database, err := db.Open()
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer database.Close()

		note, err := database.GetNote(noteTitle)
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
}
