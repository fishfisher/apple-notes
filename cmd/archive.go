package cmd

import (
	"fmt"
	"time"

	"github.com/fishfisher/apple-notes/internal/applescript"
	"github.com/fishfisher/apple-notes/internal/db"
	"github.com/spf13/cobra"
)

var (
	archiveFolder     string
	archiveOlderThan  int
	archiveTargetName string
)

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive old notes",
	Long:  `Move old notes to an archive folder based on modification date.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		database, err := db.Open()
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer database.Close()

		notes, err := database.ListNotes(archiveFolder)
		if err != nil {
			return fmt.Errorf("failed to list notes: %w", err)
		}

		// Calculate cutoff date
		cutoffMonths := archiveOlderThan
		cutoffDate := time.Now().AddDate(0, -cutoffMonths, 0)

		var toArchive []db.Note
		for _, note := range notes {
			if note.Modified.Before(cutoffDate) {
				toArchive = append(toArchive, note)
			}
		}

		if len(toArchive) == 0 {
			fmt.Printf("No notes older than %d months found\n", cutoffMonths)
			return nil
		}

		fmt.Printf("Found %d notes older than %d months\n", len(toArchive), cutoffMonths)
		fmt.Printf("Move to '%s' folder? (y/N): ", archiveTargetName)

		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Archive cancelled")
			return nil
		}

		// Move notes
		moved := 0
		for _, note := range toArchive {
			if err := applescript.MoveNote(note.Title, archiveTargetName); err != nil {
				fmt.Printf("Warning: failed to move '%s': %v\n", note.Title, err)
				continue
			}
			moved++
		}

		fmt.Printf("Successfully archived %d notes to '%s'\n", moved, archiveTargetName)
		return nil
	},
}

func init() {
	archiveCmd.Flags().StringVarP(&archiveFolder, "folder", "f", "", "Source folder to archive from (empty = all)")
	archiveCmd.Flags().IntVarP(&archiveOlderThan, "older-than", "o", 6, "Archive notes older than N months")
	archiveCmd.Flags().StringVarP(&archiveTargetName, "to", "t", "Archive", "Target folder name")
}
