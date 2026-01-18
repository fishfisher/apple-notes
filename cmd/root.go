package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version string
	commit  string
	date    string
)

var rootCmd = &cobra.Command{
	Use:   "apple-notes",
	Short: "CLI for Apple Notes",
	Long: `apple-notes is a command-line interface for Apple Notes.
It uses SQLite for fast read operations and AppleScript for write operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func SetVersionInfo(v, c, d string) {
	version = v
	commit = c
	date = d
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(foldersCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(moveCmd)
	rootCmd.AddCommand(versionCmd)
}
