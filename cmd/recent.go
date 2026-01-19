package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fishfisher/apple-notes/internal/db"
	"github.com/spf13/cobra"
)

var (
	recentToday bool
	recentWeek  bool
	recentLimit int
)

var recentCmd = &cobra.Command{
	Use:   "recent",
	Short: "Show recently modified notes",
	Long:  `Display recently modified notes. Use --today, --week, or specify custom days.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		database, err := db.Open()
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer database.Close()

		days := 30 // default
		if recentToday {
			days = 1
		} else if recentWeek {
			days = 7
		}

		limit := recentLimit
		if limit == 0 {
			limit = 20
		}

		notes, err := database.GetRecentNotes(days, limit)
		if err != nil {
			return fmt.Errorf("failed to get recent notes: %w", err)
		}

		if len(notes) == 0 {
			fmt.Println("No recent notes found")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTITLE\tFOLDER\tMODIFIED\tSNIPPET")
		for _, note := range notes {
			snippet := note.Snippet
			if len(snippet) > 60 {
				snippet = snippet[:60] + "..."
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				note.ID,
				note.Title,
				note.Folder,
				note.Modified.Format("2006-01-02 15:04"),
				snippet,
			)
		}
		w.Flush()

		periodStr := fmt.Sprintf("last %d days", days)
		if recentToday {
			periodStr = "today"
		} else if recentWeek {
			periodStr = "this week"
		}
		fmt.Printf("\nShowing %d notes from %s\n", len(notes), periodStr)
		return nil
	},
}

func init() {
	recentCmd.Flags().BoolVar(&recentToday, "today", false, "Show notes modified today")
	recentCmd.Flags().BoolVar(&recentWeek, "week", false, "Show notes modified this week")
	recentCmd.Flags().IntVarP(&recentLimit, "limit", "l", 20, "Maximum number of notes to show")
}
