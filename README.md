# apple-notes

A blazing-fast command-line interface for Apple Notes. Directly queries the SQLite database instead of using AppleScript, making it orders of magnitude faster for large note collections.

## Why apple-notes?

The popular [`memo`](https://github.com/antoniorodr/memo) CLI uses AppleScript to access Apple Notes, which becomes painfully slow with large note collections (thousands of notes). `apple-notes` solves this by directly querying the Notes SQLite database, providing instant results even with massive note libraries.

### Performance Comparison

- **memo search**: 30+ seconds (hangs with large collections)
- **apple-notes search**: <100ms

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

- `search [term]` - Search notes by title or content
- `list` - List all notes (supports `--folder` and `--limit` flags)
- `show [note-id-or-title]` - Show a specific note
- `folders` - List all folders with note counts
- `export` - Export notes to JSON or text format
- `version` - Print version information

## How It Works

`apple-notes` directly queries the Apple Notes SQLite database located at:
```
~/Library/Group Containers/group.com.apple.notes/NoteStore.sqlite
```

The database is opened in read-only mode to avoid any locking issues or conflicts with the Notes app.

## Limitations

- **Read-only**: This tool can only read notes, not create or modify them. Use the Notes app for editing.
- **macOS only**: Apple Notes database is only available on macOS.
- **Snippet only**: The full note body with formatting is stored in a complex format. This tool currently shows the plain text snippet.

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
