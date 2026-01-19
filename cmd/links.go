package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fishfisher/apple-notes/internal/db"
	"github.com/spf13/cobra"
)

var (
	linksAll     bool
	linksByTitle bool
)

var linksCmd = &cobra.Command{
	Use:   "links [note-id-or-title]",
	Short: "Extract URLs from notes",
	Long:  `Extract all URLs from a specific note or find all notes containing URLs.

By default, numeric input is treated as a note ID. Use --by-title to search by title instead.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		database, err := db.Open()
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer database.Close()

		if linksAll {
			// Find all notes with links
			notes, err := database.FindNotesWithLinks()
			if err != nil {
				return fmt.Errorf("failed to find notes with links: %w", err)
			}

			if len(notes) == 0 {
				fmt.Println("No notes with links found")
				return nil
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tTITLE\tFOLDER\tMODIFIED")
			for _, note := range notes {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
					note.ID,
					note.Title,
					note.Folder,
					note.Modified.Format("2006-01-02 15:04"),
				)
			}
			w.Flush()

			fmt.Printf("\nFound %d notes with links\n", len(notes))
			return nil
		}

		// Extract links from specific note
		if len(args) == 0 {
			return fmt.Errorf("note ID or title required (or use --all flag)")
		}

		identifier := args[0]
		var note *db.Note
		if linksByTitle {
			note, err = database.GetNoteByTitle(identifier)
		} else {
			note, err = database.GetNote(identifier)
		}
		if err != nil {
			return fmt.Errorf("note not found: %w", err)
		}

		urls, err := database.ExtractLinks(note.ID)
		if err != nil {
			return fmt.Errorf("failed to extract links: %w", err)
		}

		if len(urls) == 0 {
			fmt.Printf("No links found in note '%s'\n", note.Title)
			return nil
		}

		fmt.Printf("Links in '%s':\n", note.Title)
		for i, url := range urls {
			fmt.Printf("%d. %s\n", i+1, url)
		}

		return nil
	},
}

func init() {
	linksCmd.Flags().BoolVarP(&linksAll, "all", "a", false, "Find all notes containing links")
	linksCmd.Flags().BoolVarP(&linksByTitle, "by-title", "t", false, "Search by title (even if input is numeric)")
}
