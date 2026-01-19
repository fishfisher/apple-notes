package cmd

import (
	"fmt"

	"github.com/fishfisher/apple-notes/internal/db"
	"github.com/spf13/cobra"
)

var duplicatesCmd = &cobra.Command{
	Use:   "duplicates",
	Short: "Find duplicate notes",
	Long:  `Find notes with identical titles.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		database, err := db.Open()
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer database.Close()

		duplicates, err := database.FindDuplicates()
		if err != nil {
			return fmt.Errorf("failed to find duplicates: %w", err)
		}

		if len(duplicates) == 0 {
			fmt.Println("No duplicate notes found")
			return nil
		}

		fmt.Printf("Found %d sets of duplicate notes:\n\n", len(duplicates))

		for i, group := range duplicates {
			fmt.Printf("%d. Title: %s (%d copies)\n", i+1, group[0].Title, len(group))
			for j, note := range group {
				fmt.Printf("   %c. ID: %s, Folder: %s, Modified: %s\n",
					'a'+j,
					note.ID,
					note.Folder,
					note.Modified.Format("2006-01-02 15:04"),
				)
			}
			fmt.Println()
		}

		return nil
	},
}
