package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fishfisher/apple-notes/internal/applescript"
	"github.com/fishfisher/apple-notes/internal/db"
	"github.com/spf13/cobra"
)

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "Manage tags (hashtags)",
	Long:  `List, search, and manage hashtags in notes.`,
}

var tagsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tags with counts",
	RunE: func(cmd *cobra.Command, args []string) error {
		database, err := db.Open()
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer database.Close()

		tags, err := database.ExtractTags()
		if err != nil {
			return fmt.Errorf("failed to extract tags: %w", err)
		}

		if len(tags) == 0 {
			fmt.Println("No tags found")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TAG\tCOUNT")
		for _, tag := range tags {
			fmt.Fprintf(w, "%s\t%d\n", tag.Name, tag.Count)
		}
		w.Flush()

		fmt.Printf("\nTotal: %d unique tags\n", len(tags))
		return nil
	},
}

var tagsSearchCmd = &cobra.Command{
	Use:   "search [tag]",
	Short: "Search notes by tag",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tag := args[0]

		database, err := db.Open()
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer database.Close()

		notes, err := database.SearchByTag(tag)
		if err != nil {
			return fmt.Errorf("failed to search by tag: %w", err)
		}

		if len(notes) == 0 {
			fmt.Printf("No notes found with tag '%s'\n", tag)
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTITLE\tFOLDER\tMODIFIED")
		for _, note := range notes {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				note.ID,
				note.Title,
				note.Folder,
				note.Modified.Format("2006-01-02 15:04"),
			)
		}
		w.Flush()

		fmt.Printf("\nFound %d notes with tag '%s'\n", len(notes), tag)
		return nil
	},
}

var tagsAddCmd = &cobra.Command{
	Use:   "add [note-id] [tag]",
	Short: "Add a tag to a note",
	Long:  `Add a hashtag to a note by ID.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		noteID := args[0]
		tag := args[1]

		// Verify note exists
		database, err := db.Open()
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer database.Close()

		note, err := database.GetNote(noteID)
		if err != nil {
			return fmt.Errorf("note not found: %w", err)
		}

		fmt.Printf("Adding tag '%s' to note '%s'...\n", tag, note.Title)
		if err := applescript.AddTagToNote(note.Title, tag); err != nil {
			return fmt.Errorf("failed to add tag: %w", err)
		}

		fmt.Println("Tag added successfully")
		return nil
	},
}

func init() {
	tagsCmd.AddCommand(tagsListCmd)
	tagsCmd.AddCommand(tagsSearchCmd)
	tagsCmd.AddCommand(tagsAddCmd)
}
