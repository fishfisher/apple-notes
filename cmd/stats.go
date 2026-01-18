package cmd

import (
	"fmt"

	"github.com/fishfisher/apple-notes/internal/db"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show statistics about your notes collection",
	Long:  `Display statistics including note counts, top tags, and more.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		database, err := db.Open()
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer database.Close()

		stats, err := database.GetStats()
		if err != nil {
			return fmt.Errorf("failed to get stats: %w", err)
		}

		fmt.Println("=== Apple Notes Statistics ===\n")
		fmt.Printf("Total notes:           %d\n", stats.TotalNotes)
		fmt.Printf("Total folders:         %d\n", stats.TotalFolders)
		fmt.Printf("Modified this week:    %d\n", stats.NotesThisWeek)
		fmt.Printf("Modified this month:   %d\n", stats.NotesThisMonth)

		if stats.LargestNote.Title != "" {
			fmt.Printf("\nLargest note:        %s\n", stats.LargestNote.Title)
			fmt.Printf("  Folder:            %s\n", stats.LargestNote.Folder)
			fmt.Printf("  Size:              %d characters\n", stats.TotalCharacters)
		}

		if len(stats.TopTags) > 0 {
			fmt.Println("\nTop tags:")
			for i, tag := range stats.TopTags {
				if i >= 10 {
					break
				}
				fmt.Printf("  %s (%d)\n", tag.Name, tag.Count)
			}
		}

		return nil
	},
}
