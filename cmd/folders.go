package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fishfisher/apple-notes/internal/db"
	"github.com/spf13/cobra"
)

var foldersCmd = &cobra.Command{
	Use:   "folders",
	Short: "List all folders",
	Long:  `List all note folders with note counts.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		database, err := db.Open()
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer database.Close()

		folders, err := database.ListFolders()
		if err != nil {
			return fmt.Errorf("failed to list folders: %w", err)
		}

		if len(folders) == 0 {
			fmt.Println("No folders found")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "FOLDER\tNOTES")
		totalNotes := 0
		for _, folder := range folders {
			fmt.Fprintf(w, "%s\t%d\n", folder.Name, folder.Count)
			totalNotes += folder.Count
		}
		w.Flush()

		fmt.Printf("\nTotal: %d folders, %d notes\n", len(folders), totalNotes)
		return nil
	},
}
