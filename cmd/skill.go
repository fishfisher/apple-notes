package cmd

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	skillFS      fs.FS
	skillDirName string
	skillName    string
)

// SetSkillFS receives the embedded filesystem and metadata from main.
func SetSkillFS(fsys fs.FS, dirName, name string) {
	skillFS = fsys
	skillDirName = dirName
	skillName = name
}

var installSkillCmd = &cobra.Command{
	Use:   "install-skill",
	Short: "Install the Claude Code skill for this CLI",
	Long:  "Copies the embedded skill files to ~/.claude/skills/ so Claude Code can discover them globally.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if skillFS == nil {
			return fmt.Errorf("no skill files embedded")
		}

		baseDir, _ := cmd.Flags().GetString("path")
		force, _ := cmd.Flags().GetBool("force")

		if baseDir == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("finding home directory: %w", err)
			}
			baseDir = filepath.Join(home, ".claude", "skills")
		}
		destDir := filepath.Join(baseDir, skillName)

		// Collect all files from the embedded FS
		type fileEntry struct {
			relPath string // path relative to destDir
			embPath string // path inside embed.FS
		}
		var files []fileEntry

		err := fs.WalkDir(skillFS, skillDirName, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			rel, _ := filepath.Rel(skillDirName, path)
			files = append(files, fileEntry{relPath: rel, embPath: path})
			return nil
		})
		if err != nil {
			return fmt.Errorf("reading embedded files: %w", err)
		}

		// Check which files already exist
		var existing []string
		for _, f := range files {
			dest := filepath.Join(destDir, f.relPath)
			if _, err := os.Stat(dest); err == nil {
				existing = append(existing, f.relPath)
			}
		}

		// Show file list
		bold := color.New(color.Bold)
		bold.Printf("Installing to %s\n", destDir)
		for _, f := range files {
			marker := ""
			for _, e := range existing {
				if e == f.relPath {
					marker = color.YellowString(" (exists)")
					break
				}
			}
			fmt.Printf("  %s%s\n", f.relPath, marker)
		}

		// Prompt if existing files and not --force
		if len(existing) > 0 && !force {
			fmt.Printf("\n%s Overwrite %d existing file(s)? [y/N] ", color.YellowString("?"), len(existing))
			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')
			if strings.TrimSpace(strings.ToLower(answer)) != "y" {
				fmt.Println("Aborted.")
				return nil
			}
		}

		// Create directories and copy files
		for _, f := range files {
			dest := filepath.Join(destDir, f.relPath)
			if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
				return fmt.Errorf("creating directory: %w", err)
			}
			data, err := fs.ReadFile(skillFS, f.embPath)
			if err != nil {
				return fmt.Errorf("reading %s: %w", f.relPath, err)
			}
			if err := os.WriteFile(dest, data, 0644); err != nil {
				return fmt.Errorf("writing %s: %w", f.relPath, err)
			}
		}

		color.Green("Installed %d file(s) to %s", len(files), destDir)
		return nil
	},
}

func init() {
	installSkillCmd.Flags().StringP("path", "p", "", "Parent directory (default: ~/.claude/skills/)")
	installSkillCmd.Flags().BoolP("force", "f", false, "Overwrite existing files without prompting")
	rootCmd.AddCommand(installSkillCmd)
}
