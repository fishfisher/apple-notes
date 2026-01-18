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
	Short: "Fast CLI for Apple Notes using SQLite",
	Long: `apple-notes is a blazing-fast command-line interface for Apple Notes.
It directly queries the SQLite database instead of using AppleScript,
making it orders of magnitude faster for large note collections.`,
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
	rootCmd.AddCommand(versionCmd)
}
