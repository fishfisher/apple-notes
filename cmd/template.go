package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fishfisher/apple-notes/internal/applescript"
	"github.com/spf13/cobra"
)

type Template struct {
	Name string `json:"name"`
	Body string `json:"body"`
}

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage note templates",
	Long:  `Create, list, and use note templates.`,
}

var templateCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new template",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		templateName := args[0]
		body, _ := cmd.Flags().GetString("body")

		if body == "" {
			return fmt.Errorf("--body flag is required")
		}

		template := Template{
			Name: templateName,
			Body: body,
		}

		if err := saveTemplate(template); err != nil {
			return fmt.Errorf("failed to save template: %w", err)
		}

		fmt.Printf("Template '%s' created successfully\n", templateName)
		return nil
	},
}

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		templates, err := loadTemplates()
		if err != nil {
			return fmt.Errorf("failed to load templates: %w", err)
		}

		if len(templates) == 0 {
			fmt.Println("No templates found")
			return nil
		}

		fmt.Println("Available templates:")
		for _, tmpl := range templates {
			fmt.Printf("  - %s\n", tmpl.Name)
		}

		return nil
	},
}

var templateUseCmd = &cobra.Command{
	Use:   "use [template-name] [note-title]",
	Short: "Create a note from a template",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		templateName := args[0]
		noteTitle := args[1]
		folder, _ := cmd.Flags().GetString("folder")

		templates, err := loadTemplates()
		if err != nil {
			return fmt.Errorf("failed to load templates: %w", err)
		}

		var template *Template
		for _, tmpl := range templates {
			if tmpl.Name == templateName {
				template = &tmpl
				break
			}
		}

		if template == nil {
			return fmt.Errorf("template '%s' not found", templateName)
		}

		if folder == "" {
			folder = "Notes"
		}

		fmt.Printf("Creating note '%s' from template '%s'...\n", noteTitle, templateName)
		if err := applescript.AddNote(noteTitle, template.Body, folder); err != nil {
			return fmt.Errorf("failed to create note: %w", err)
		}

		fmt.Println("Note created successfully")
		return nil
	},
}

func getTemplatesPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".apple-notes-templates.json"), nil
}

func loadTemplates() ([]Template, error) {
	path, err := getTemplatesPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []Template{}, nil
	}
	if err != nil {
		return nil, err
	}

	var templates []Template
	if err := json.Unmarshal(data, &templates); err != nil {
		return nil, err
	}

	return templates, nil
}

func saveTemplate(template Template) error {
	templates, err := loadTemplates()
	if err != nil {
		return err
	}

	// Check if template already exists
	for i, tmpl := range templates {
		if tmpl.Name == template.Name {
			templates[i] = template
			goto save
		}
	}

	templates = append(templates, template)

save:
	data, err := json.MarshalIndent(templates, "", "  ")
	if err != nil {
		return err
	}

	path, err := getTemplatesPath()
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func init() {
	templateCreateCmd.Flags().String("body", "", "Template body")
	templateCreateCmd.MarkFlagRequired("body")

	templateUseCmd.Flags().StringP("folder", "f", "Notes", "Folder to create note in")

	templateCmd.AddCommand(templateCreateCmd)
	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateUseCmd)
}
