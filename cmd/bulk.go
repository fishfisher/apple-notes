package cmd

import (
	"fmt"

	"github.com/fishfisher/apple-notes/internal/applescript"
	"github.com/spf13/cobra"
)

var bulkCmd = &cobra.Command{
	Use:   "bulk",
	Short: "Bulk operations on notes",
	Long:  `Perform bulk operations like moving all notes from one folder to another.`,
}

var bulkMoveCmd = &cobra.Command{
	Use:   "move --from [source-folder] --to [target-folder]",
	Short: "Move all notes from one folder to another",
	RunE: func(cmd *cobra.Command, args []string) error {
		sourceFolder, _ := cmd.Flags().GetString("from")
		targetFolder, _ := cmd.Flags().GetString("to")

		if sourceFolder == "" || targetFolder == "" {
			return fmt.Errorf("both --from and --to flags are required")
		}

		fmt.Printf("Move all notes from '%s' to '%s'? (y/N): ", sourceFolder, targetFolder)
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Bulk move cancelled")
			return nil
		}

		fmt.Printf("Moving notes from '%s' to '%s'...\n", sourceFolder, targetFolder)
		count, err := applescript.BulkMoveNotes(sourceFolder, targetFolder)
		if err != nil {
			return fmt.Errorf("failed to move notes: %w", err)
		}

		fmt.Printf("Successfully moved %d notes\n", count)
		return nil
	},
}

func init() {
	bulkMoveCmd.Flags().String("from", "", "Source folder")
	bulkMoveCmd.Flags().String("to", "", "Target folder")
	bulkMoveCmd.MarkFlagRequired("from")
	bulkMoveCmd.MarkFlagRequired("to")

	bulkCmd.AddCommand(bulkMoveCmd)
}
