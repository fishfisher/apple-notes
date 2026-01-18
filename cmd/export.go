package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fishfisher/apple-notes/internal/db"
	"github.com/spf13/cobra"
)

var (
	exportOutput string
	exportFormat string
	exportFolder string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export notes to a file",
	Long:  `Export notes to JSON or text format.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		database, err := db.Open()
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer database.Close()

		notes, err := database.ListNotes(exportFolder)
		if err != nil {
			return fmt.Errorf("failed to list notes: %w", err)
		}

		if len(notes) == 0 {
			fmt.Println("No notes to export")
			return nil
		}

		// Default output path
		if exportOutput == "" {
			home, _ := os.UserHomeDir()
			exportOutput = filepath.Join(home, "Desktop", "apple-notes-export.json")
			if exportFormat == "txt" {
				exportOutput = filepath.Join(home, "Desktop", "apple-notes-export.txt")
			}
		}

		var data []byte
		switch exportFormat {
		case "json":
			data, err = json.MarshalIndent(notes, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
		case "txt":
			var text string
			for _, note := range notes {
				text += fmt.Sprintf("Title: %s\n", note.Title)
				text += fmt.Sprintf("Folder: %s\n", note.Folder)
				text += fmt.Sprintf("Modified: %s\n", note.Modified.Format("2006-01-02 15:04:05"))
				text += fmt.Sprintf("\n%s\n", note.Snippet)
				text += "\n" + string(make([]byte, 80)) + "\n\n"
			}
			data = []byte(text)
		default:
			return fmt.Errorf("unsupported format: %s (use json or txt)", exportFormat)
		}

		if err := os.WriteFile(exportOutput, data, 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		fmt.Printf("Exported %d notes to %s\n", len(notes), exportOutput)
		return nil
	},
}

func init() {
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "Output file path (default: ~/Desktop/apple-notes-export.json)")
	exportCmd.Flags().StringVarP(&exportFormat, "format", "t", "json", "Export format: json or txt")
	exportCmd.Flags().StringVarP(&exportFolder, "folder", "f", "", "Export only notes from this folder")
}
