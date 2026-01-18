package applescript

import (
	"fmt"
	"os/exec"
	"strings"
)

// execAppleScript executes an AppleScript and returns the output
func execAppleScript(script string) (string, error) {
	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("AppleScript error: %w\nOutput: %s", err, string(output))
	}
	return strings.TrimSpace(string(output)), nil
}

// AddNote creates a new note with the given title and body in the specified folder
func AddNote(title, body, folder string) error {
	script := fmt.Sprintf(`
		tell application "Notes"
			tell folder "%s"
				make new note with properties {name:"%s", body:"%s"}
			end tell
		end tell
	`, escapeQuotes(folder), escapeQuotes(title), escapeQuotes(body))

	_, err := execAppleScript(script)
	return err
}

// EditNote updates an existing note's title and/or body
func EditNote(noteTitle, newTitle, newBody string) error {
	script := fmt.Sprintf(`
		tell application "Notes"
			set theNote to first note whose name is "%s"
			set body of theNote to "%s"
			set name of theNote to "%s"
		end tell
	`, escapeQuotes(noteTitle), escapeQuotes(newBody), escapeQuotes(newTitle))

	_, err := execAppleScript(script)
	return err
}

// DeleteNote deletes a note by title
func DeleteNote(noteTitle string) error {
	script := fmt.Sprintf(`
		tell application "Notes"
			delete (first note whose name is "%s")
		end tell
	`, escapeQuotes(noteTitle))

	_, err := execAppleScript(script)
	return err
}

// MoveNote moves a note to a different folder
func MoveNote(noteTitle, targetFolder string) error {
	script := fmt.Sprintf(`
		tell application "Notes"
			set theNote to first note whose name is "%s"
			move theNote to folder "%s"
		end tell
	`, escapeQuotes(noteTitle), escapeQuotes(targetFolder))

	_, err := execAppleScript(script)
	return err
}

// ListFolderNames returns all folder names
func ListFolderNames() ([]string, error) {
	script := `
		tell application "Notes"
			set folderList to {}
			repeat with f in folders
				set end of folderList to name of f
			end repeat
			return folderList
		end tell
	`

	output, err := execAppleScript(script)
	if err != nil {
		return nil, err
	}

	// Parse comma-separated list
	if output == "" {
		return []string{}, nil
	}

	folders := strings.Split(output, ", ")
	return folders, nil
}

// AppendNote appends content to an existing note
func AppendNote(noteTitle, content string) error {
	script := fmt.Sprintf(`
		tell application "Notes"
			set theNote to first note whose name is "%s"
			set body of theNote to (body of theNote & "\n%s")
		end tell
	`, escapeQuotes(noteTitle), escapeQuotes(content))

	_, err := execAppleScript(script)
	return err
}

// AddTagToNote adds a hashtag to a note (appends it to the body)
func AddTagToNote(noteTitle, tag string) error {
	// Ensure tag starts with #
	if !strings.HasPrefix(tag, "#") {
		tag = "#" + tag
	}

	script := fmt.Sprintf(`
		tell application "Notes"
			set theNote to first note whose name is "%s"
			set body of theNote to (body of theNote & " %s")
		end tell
	`, escapeQuotes(noteTitle), escapeQuotes(tag))

	_, err := execAppleScript(script)
	return err
}

// BulkMoveNotes moves all notes from a source folder to a target folder
func BulkMoveNotes(sourceFolder, targetFolder string) (int, error) {
	script := fmt.Sprintf(`
		tell application "Notes"
			set movedCount to 0
			set sourceNotes to notes of folder "%s"
			repeat with n in sourceNotes
				move n to folder "%s"
				set movedCount to movedCount + 1
			end repeat
			return movedCount
		end tell
	`, escapeQuotes(sourceFolder), escapeQuotes(targetFolder))

	output, err := execAppleScript(script)
	if err != nil {
		return 0, err
	}

	var count int
	fmt.Sscanf(output, "%d", &count)
	return count, nil
}

// CreateFolder creates a new folder
func CreateFolder(folderName string) error {
	script := fmt.Sprintf(`
		tell application "Notes"
			make new folder with properties {name:"%s"}
		end tell
	`, escapeQuotes(folderName))

	_, err := execAppleScript(script)
	return err
}

// escapeQuotes escapes double quotes for AppleScript
func escapeQuotes(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}
