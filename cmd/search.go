package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fishfisher/apple-notes/internal/db"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [term]",
	Short: "Search notes by title or content",
	Long:  `Search for notes containing the specified term in their title or snippet.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		searchTerm := args[0]

		database, err := db.Open()
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer database.Close()

		notes, err := database.SearchNotes(searchTerm)
		if err != nil {
			return fmt.Errorf("failed to search notes: %w", err)
		}

		if len(notes) == 0 {
			fmt.Printf("No notes found matching '%s'\n", searchTerm)
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TITLE\tFOLDER\tMODIFIED\tSNIPPET")
		for _, note := range notes {
			snippet := note.Snippet
			if len(snippet) > 60 {
				snippet = snippet[:60] + "..."
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				note.Title,
				note.Folder,
				note.Modified.Format("2006-01-02 15:04"),
				snippet,
			)
		}
		w.Flush()

		fmt.Printf("\nFound %d notes matching '%s'\n", len(notes), searchTerm)
		return nil
	},
}
