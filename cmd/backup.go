package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fishfisher/apple-notes/internal/applescript"
	"github.com/fishfisher/apple-notes/internal/db"
	"github.com/spf13/cobra"
)

type Backup struct {
	Timestamp time.Time `json:"timestamp"`
	Notes     []db.Note `json:"notes"`
}

var backupCmd = &cobra.Command{
	Use:   "backup [output-path]",
	Short: "Backup all notes to a JSON file",
	Long:  `Create a complete backup of all notes in JSON format.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		database, err := db.Open()
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer database.Close()

		notes, err := database.ListNotes("")
		if err != nil {
			return fmt.Errorf("failed to list notes: %w", err)
		}

		backup := Backup{
			Timestamp: time.Now(),
			Notes:     notes,
		}

		data, err := json.MarshalIndent(backup, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal backup: %w", err)
		}

		outputPath := ""
		if len(args) > 0 {
			outputPath = args[0]
		} else {
			home, _ := os.UserHomeDir()
			timestamp := time.Now().Format("2006-01-02-150405")
			outputPath = filepath.Join(home, "Desktop", fmt.Sprintf("apple-notes-backup-%s.json", timestamp))
		}

		if err := os.WriteFile(outputPath, data, 0644); err != nil {
			return fmt.Errorf("failed to write backup: %w", err)
		}

		fmt.Printf("Backed up %d notes to %s\n", len(notes), outputPath)
		return nil
	},
}

var restoreCmd = &cobra.Command{
	Use:   "restore [backup-path]",
	Short: "Restore notes from a backup file",
	Long:  `Restore notes from a JSON backup file. This will create new notes from the backup.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		backupPath := args[0]

		data, err := os.ReadFile(backupPath)
		if err != nil {
			return fmt.Errorf("failed to read backup: %w", err)
		}

		var backup Backup
		if err := json.Unmarshal(data, &backup); err != nil {
			return fmt.Errorf("failed to parse backup: %w", err)
		}

		fmt.Printf("Backup from %s contains %d notes\n", backup.Timestamp.Format("2006-01-02 15:04:05"), len(backup.Notes))
		fmt.Printf("Restore these notes? (y/N): ")

		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Restore cancelled")
			return nil
		}

		restored := 0
		for i, note := range backup.Notes {
			fmt.Printf("Restoring %d/%d: %s\n", i+1, len(backup.Notes), note.Title)
			if err := applescript.AddNote(note.Title, note.Snippet, note.Folder); err != nil {
				fmt.Printf("  Warning: failed to restore '%s': %v\n", note.Title, err)
				continue
			}
			restored++
		}

		fmt.Printf("\nSuccessfully restored %d/%d notes\n", restored, len(backup.Notes))
		return nil
	},
}
