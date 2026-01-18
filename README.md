# apple-notes

A command-line interface for Apple Notes that uses SQLite for fast read operations and AppleScript for write operations.

## Features

- **Fast search and listing**: Direct SQLite queries for instant results with large note collections
- **Full CRUD operations**: Create, read, update, and delete notes
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

## Commands

### Read Operations (SQLite-based - Fast)
- `search [term]` - Search notes by title or content
- `list` - List all notes (supports `--folder` and `--limit` flags)
- `show [note-id-or-title]` - Show a specific note
- `folders` - List all folders with note counts
- `export` - Export notes to JSON or text format

### Write Operations (AppleScript-based)
- `add [title]` - Create a new note
- `edit [note-title]` - Edit an existing note
- `delete [note-title]` - Delete a note
- `move [note-title] [folder]` - Move a note to a different folder

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
