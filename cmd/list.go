package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fishfisher/apple-notes/internal/db"
	"github.com/spf13/cobra"
)

var (
	listFolder string
	listLimit  int
	listHideID bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notes",
	Long:  `List all notes, optionally filtered by folder.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		database, err := db.Open()
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer database.Close()

		notes, err := database.ListNotes(listFolder)
		if err != nil {
			return fmt.Errorf("failed to list notes: %w", err)
		}

		if len(notes) == 0 {
			fmt.Println("No notes found")
			return nil
		}

		// Apply limit if specified
		if listLimit > 0 && len(notes) > listLimit {
			notes = notes[:listLimit]
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		if listHideID {
			fmt.Fprintln(w, "TITLE\tFOLDER\tMODIFIED\tSNIPPET")
		} else {
			fmt.Fprintln(w, "ID\tTITLE\tFOLDER\tMODIFIED\tSNIPPET")
		}
		for _, note := range notes {
			snippet := note.Snippet
			if len(snippet) > 60 {
				snippet = snippet[:60] + "..."
			}
			if listHideID {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
					note.Title,
					note.Folder,
					note.Modified.Format("2006-01-02 15:04"),
					snippet,
				)
			} else {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					note.ID,
					note.Title,
					note.Folder,
					note.Modified.Format("2006-01-02 15:04"),
					snippet,
				)
			}
		}
		w.Flush()

		fmt.Printf("\nTotal: %d notes\n", len(notes))
		return nil
	},
}

func init() {
	listCmd.Flags().StringVarP(&listFolder, "folder", "f", "", "Filter by folder name")
	listCmd.Flags().IntVarP(&listLimit, "limit", "l", 0, "Limit number of results (0 = no limit)")
	listCmd.Flags().BoolVar(&listHideID, "hide-id", false, "Hide note IDs from output")
}
