# apple-notes

A command-line interface for Apple Notes that uses SQLite for fast read operations and AppleScript for write operations.

## Features

- **Fast search and listing**: Direct SQLite queries for instant results with large note collections
- **Full CRUD operations**: Create, read, update, and delete notes
- **Quick append**: Quickly add content to existing notes
- **Tags management**: List, search, and manage hashtags
- **Recent notes**: Filter notes by modification date
- **Statistics**: View analytics about your note collection
- **Duplicates detection**: Find notes with identical titles
- **Link extraction**: Find notes with URLs
- **Templates**: Create and reuse note templates
- **Bulk operations**: Move entire folders at once
- **Archive**: Automatically archive old notes
- **Backup/Restore**: Full backup and restore functionality
- **Folder management**: Organize notes across folders
- **Export**: Export notes to JSON or text format

## Installation

### Homebrew

```bash
brew install fishfisher/tap/apple-notes
```

### From Source

```bash
go install github.com/fishfisher/apple-notes@latest
```

### Binary Download

Download the latest release from the [releases page](https://github.com/fishfisher/apple-notes/releases).

## Usage

### Search notes

```bash
apple-notes search "meeting notes"
```

### List all notes

```bash
apple-notes list
```

### List notes from a specific folder

```bash
apple-notes list --folder Work
```

### Show a specific note

```bash
apple-notes show "My Note Title"
```

### Add a new note

```bash
# Interactive mode (prompts for body)
apple-notes add "My New Note" --folder Work

# With body from flag
apple-notes add "Shopping List" --body "Milk, Eggs, Bread" --folder Personal
```

### Edit a note

```bash
# Interactive mode (prompts for new body)
apple-notes edit "My Note Title"

# With body from flag
apple-notes edit "Shopping List" --body "Milk, Eggs, Bread, Butter"

# Rename a note
apple-notes edit "Old Title" --title "New Title" --body "Updated content"

# Notes with images/attachments are protected by default
# The tool will block edits and suggest using Notes.app
# Use --force-unsafe to override (NOT RECOMMENDED - destroys rich content)
```

### Delete a note

```bash
# With confirmation prompt
apple-notes delete "My Note Title"

# Skip confirmation
apple-notes delete "My Note Title" --force
```

### Move a note to another folder

```bash
apple-notes move "My Note Title" "Work"
```

### List all folders

```bash
apple-notes folders
```

### Export notes

```bash
# Export to JSON (default)
apple-notes export

# Export to text
apple-notes export --format txt

# Export specific folder
apple-notes export --folder Work --output ~/Desktop/work-notes.json
```

### Append to a note

```bash
# Interactive mode
apple-notes append "Daily Log"

# With content flag
apple-notes append "Daily Log" --content "Meeting with Sarah at 2pm"
```

### Recent notes

```bash
# Show recently modified notes (default: last 30 days, limit 20)
apple-notes recent

# Show notes from today
apple-notes recent --today

# Show notes from this week
apple-notes recent --week

# Custom limit
apple-notes recent --limit 10
```

### Tags (hashtags)

```bash
# List all tags with counts
apple-notes tags list

# Search notes by tag
apple-notes tags search "#work"

# Add tag to a note
apple-notes tags add "My Note" "#important"
```

### Statistics

```bash
# Show collection statistics
apple-notes stats
```

### Find duplicates

```bash
# Find notes with identical titles
apple-notes duplicates
```

### Extract links

```bash
# Extract URLs from a specific note
apple-notes links "My Note"

# Find all notes containing URLs
apple-notes links --all
```

### Templates

```bash
# Create a template
apple-notes template create "meeting" --body "Date:\nAttendees:\nAgenda:\nNotes:"

# List templates
apple-notes template list

# Create note from template
apple-notes template use "meeting" "Team Standup" --folder Work
```

### Bulk operations

```bash
# Move all notes from one folder to another
apple-notes bulk move --from "Old Projects" --to "Archive"
```

### Archive old notes

```bash
# Archive notes older than 6 months (default)
apple-notes archive --folder "Work" --to "Archive"

# Archive notes older than 1 year
apple-notes archive --older-than 12 --to "Archive"
```

### Backup and restore

```bash
# Backup all notes
apple-notes backup ~/backups/notes-2026-01-18.json

# Restore from backup
apple-notes restore ~/backups/notes-2026-01-18.json
```

## Commands

### Read Operations (SQLite-based - Fast)
- `search [term]` - Search notes by title or content
- `list` - List all notes (supports `--folder` and `--limit` flags)
- `show [note-id-or-title]` - Show a specific note
- `folders` - List all folders with note counts
- `recent` - Show recently modified notes (supports `--today`, `--week`, `--limit`)
- `stats` - Display collection statistics
- `duplicates` - Find notes with identical titles
- `links [note-title]` - Extract URLs from notes (use `--all` to find all notes with links)
- `export` - Export notes to JSON or text format

### Write Operations (AppleScript-based)
- `add [title]` - Create a new note
- `edit [note-title]` - Edit an existing note
- `delete [note-title]` - Delete a note
- `move [note-title] [folder]` - Move a note to a different folder
- `append [note-title]` - Append content to an existing note

### Tags Management
- `tags list` - List all tags with counts
- `tags search [tag]` - Search notes by tag
- `tags add [note-title] [tag]` - Add a tag to a note

### Bulk Operations
- `bulk move --from [folder] --to [folder]` - Move all notes from one folder to another
- `archive` - Archive old notes (supports `--folder`, `--older-than`, `--to`)

### Templates
- `template create [name] --body [content]` - Create a new template
- `template list` - List all templates
- `template use [template] [note-title]` - Create a note from a template

### Backup & Restore
- `backup [output-path]` - Backup all notes to JSON
- `restore [backup-path]` - Restore notes from a backup file

### Other
- `version` - Print version information

## How It Works

`apple-notes` uses a hybrid approach:

**Read operations** directly query the Apple Notes SQLite database for fast results:
```
~/Library/Group Containers/group.com.apple.notes/NoteStore.sqlite
```

**Write operations** use AppleScript to safely modify notes through the Notes app API, ensuring proper sync and data integrity.

## Limitations

- **macOS only**: Apple Notes database is only available on macOS.
- **Snippet only for display**: The full note body with formatting is stored in a complex binary format. Read commands show the plain text snippet.
- **✅ Smart Edit Protection**: The tool automatically detects notes with rich content (images, attachments, tables, sketches) and **blocks edits by default** to prevent data loss. You'll be directed to use Notes.app for those notes. This protection can be overridden with `--force-unsafe`, but this is **NOT RECOMMENDED** as it will permanently destroy:
  - Images and photos
  - Attachments and files
  - Tables and formatting
  - Sketches and drawings
  - Any other rich media

- **Append is safer**: The `append` command adds text to existing notes without replacing content, so it won't damage existing attachments. Appended content will be plain text.

## Development

```bash
# Clone the repository
git clone https://github.com/fishfisher/apple-notes.git
cd apple-notes

# Install dependencies
go mod download

# Build
go build -o apple-notes .

# Run
./apple-notes --help
```

## License

MIT

## Credits

Inspired by [memo](https://github.com/antoniorodr/memo) and [homeyctl](https://github.com/langtind/homeyctl).
