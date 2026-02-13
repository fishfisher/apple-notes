---
name: apple-notes
description: Manage Apple Notes from the command line using apple-notes CLI v0.4.0+. Use when asked to search, create, edit, delete, or organize Apple Notes, manage folders, work with tags, export notes, create templates, or perform bulk operations. All note operations use ID-based targeting for reliable, unambiguous control. Covers note CRUD operations, folder management, tags, recent notes, statistics, duplicates detection, link extraction, templates, bulk moves, archiving, and backup/restore.
metadata: {"clawdbot":{"emoji":"📝","requires":{"bins":["apple-notes"]},"install":[{"id":"brew","kind":"brew","formula":"fishfisher/tap/apple-notes","bins":["apple-notes"],"label":"Install apple-notes (brew)"}]}}
---

# Apple Notes

Manage your Apple Notes from the command line using the `apple-notes` CLI (v0.4.0+). This skill enables fast search, note creation, editing, organization, and bulk operations through a hybrid SQLite/AppleScript approach.

**Key Features:**
- **ID-only targeting**: All note operations use note IDs for reliable, unambiguous control
- **Full content display**: Commands show complete note content
- **Fast search**: Instantly find notes and get their IDs for further operations

## Setup & Configuration

### Installation

```bash
brew install fishfisher/tap/apple-notes
```

Verify installation:
```bash
apple-notes version
```

No additional configuration needed - the tool automatically accesses your Apple Notes database.

## Quick Start

```bash
# List all notes (IDs displayed for targeting)
apple-notes list

# Search for notes to get their IDs
apple-notes search "meeting"

# Show a specific note by ID
apple-notes show 4318

# Create a new note
apple-notes add "Shopping List" --body "Milk, Eggs, Bread" --folder Personal

# Edit a note by ID
apple-notes edit 4318 --body "New content"

# Append to existing note by ID
apple-notes append 4318 --content "Meeting at 2pm"

# Show recent notes
apple-notes recent --today

# List all folders
apple-notes folders

# Export notes
apple-notes export --folder Work
```

**Note Targeting:** All commands that operate on specific notes require note IDs. Use `apple-notes list` or `apple-notes search "keyword"` to find note IDs.

## How It Works

`apple-notes` uses a hybrid approach:

- **Read operations** (search, list, show) directly query the Apple Notes SQLite database for instant results with full content display
- **Write operations** (add, edit, delete) use AppleScript to safely modify notes through the Notes app, ensuring proper sync
- **ID-based targeting** provides reliable note identification - all note operations require IDs obtained from list/search commands

Database location: `~/Library/Group Containers/group.com.apple.notes/NoteStore.sqlite`

## Core Capabilities

### 1. Search & List

Find and browse notes quickly to get their IDs for use in other commands.

```bash
# Search notes by title or content (displays IDs)
apple-notes search "meeting notes"
apple-notes search "project plan"

# List all notes (IDs shown by default)
apple-notes list

# Hide IDs for cleaner output
apple-notes list --hide-id

# List notes from specific folder
apple-notes list --folder Work

# Limit results
apple-notes list --limit 10

# Show full note content by ID
apple-notes show 4318
```

### 2. Create Notes

Add new notes with content:

```bash
# Interactive mode (prompts for body)
apple-notes add "My New Note"

# With body from flag
apple-notes add "Shopping List" --body "Milk, Eggs, Bread"

# Specify folder
apple-notes add "Meeting Notes" --body "Agenda items..." --folder Work

# Quick note without body
apple-notes add "Empty Note" --folder Ideas
```

### 3. Edit Notes

Modify existing notes by ID.

```bash
# Edit note content by ID
apple-notes edit 4318 --body "New content"

# Interactive mode (prompts for new body)
apple-notes edit 4318

# Rename a note
apple-notes edit 4318 --title "New Title" --body "Updated content"
```

**Important:** Notes with images, attachments, tables, or sketches are protected by default. The tool will block edits and suggest using Notes.app. Use `--force-unsafe` to override (NOT RECOMMENDED - destroys rich content).

### 4. Append to Notes

Add content to existing notes without replacing:

```bash
# Append content by ID
apple-notes append 4318 --content "Meeting with Sarah at 2pm"

# Interactive mode (prompts for content)
apple-notes append 4318

# Append with timestamp
apple-notes append 4318 --content "$(date): Completed feature X"
```

**Note:** Append is safer than edit for notes with attachments - it won't damage existing rich content.

### 5. Delete Notes

Remove notes permanently by ID.

```bash
# Delete note by ID (with confirmation prompt)
apple-notes delete 4318

# Skip confirmation
apple-notes delete 4318 --force
```

### 6. Move Notes

Organize notes across folders by ID:

```bash
# Move note to folder
apple-notes move 4318 "Work"

# Move to new folder (creates folder if it doesn't exist)
apple-notes move 4318 "New Projects"
```

### 7. Folder Management

Manage and organize folders:

```bash
# List all folders with note counts
apple-notes folders

# List notes in specific folder
apple-notes list --folder Work

# Export entire folder
apple-notes export --folder Work --output ~/Desktop/work-notes.json
```

### 8. Recent Notes

Filter notes by modification date:

```bash
# Show recently modified notes (default: last 30 days, limit 20)
apple-notes recent

# Show notes from today
apple-notes recent --today

# Show notes from this week
apple-notes recent --week

# Custom limit
apple-notes recent --limit 10

# Combine flags
apple-notes recent --week --limit 5
```

### 9. Tags Management

Work with hashtags in notes:

```bash
# List all tags with counts
apple-notes tags list

# Search notes by tag
apple-notes tags search "#work"
apple-notes tags search "#important"

# Add tag to a note by ID
apple-notes tags add 4318 "#important"
apple-notes tags add 4318 "#work"
```

### 10. Statistics

View analytics about your note collection:

```bash
# Show collection statistics
apple-notes stats
```

Displays: total notes, notes per folder, recent activity, tag counts, etc.

### 11. Find Duplicates

Identify notes with identical titles:

```bash
# Find notes with duplicate titles
apple-notes duplicates
```

Useful for cleaning up your note collection.

### 12. Extract Links

Find and extract URLs from notes:

```bash
# Extract URLs from a specific note by ID
apple-notes links 4318

# Find all notes containing URLs
apple-notes links --all
```

### 13. Templates

Create and reuse note templates:

```bash
# Create a template
apple-notes template create "meeting" --body "Date:\nAttendees:\nAgenda:\nNotes:"

# List all templates
apple-notes template list

# Create note from template
apple-notes template use "meeting" "Team Standup" --folder Work

# Example: Daily log template
apple-notes template create "daily-log" --body "# $(date +%Y-%m-%d)\n\nTasks:\n-\n\nNotes:\n-"
```

### 14. Bulk Operations

Move entire folders at once:

```bash
# Move all notes from one folder to another
apple-notes bulk move --from "Old Projects" --to "Archive"

# Clean up by moving completed items
apple-notes bulk move --from "Active Tasks" --to "Completed"
```

### 15. Archive Old Notes

Automatically archive notes based on age:

```bash
# Archive notes older than 6 months (default)
apple-notes archive --folder "Work" --to "Archive"

# Archive notes older than 1 year
apple-notes archive --older-than 12 --to "Archive"

# Archive from specific folder
apple-notes archive --folder "Projects 2025" --older-than 6 --to "Archive/2025"
```

### 16. Backup & Restore

Full backup and restore functionality:

```bash
# Backup all notes to JSON
apple-notes backup ~/backups/notes-2026-01-18.json

# Backup with timestamp
apple-notes backup ~/backups/notes-$(date +%Y-%m-%d).json

# Restore from backup
apple-notes restore ~/backups/notes-2026-01-18.json
```

**Important:** Backups preserve note metadata but may not include rich content (images, attachments).

### 17. Export

Export notes to various formats:

```bash
# Export to JSON (default)
apple-notes export

# Export to text
apple-notes export --format txt

# Export specific folder
apple-notes export --folder Work --output ~/Desktop/work-notes.json

# Export with custom output path
apple-notes export --format txt --output ~/Documents/all-notes.txt
```

## Typical Workflow

All note operations follow a consistent pattern:

```bash
# 1. Find the note ID using search or list
apple-notes search "meeting"
# Output: [4318] Team Meeting Notes

# 2. Use the ID for all operations
apple-notes show 4318                           # View it
apple-notes edit 4318 --body "Updated content"  # Edit it
apple-notes append 4318 --content "New info"    # Append to it
apple-notes delete 4318 --force                 # Delete it
```

**Why IDs:**
- **Unique**: Each note has a permanent, unique ID
- **Fast**: ID lookups are instant
- **Reliable**: No ambiguity, even if multiple notes share similar titles

## Common Workflows

### Daily Journaling

```bash
# Create daily template
apple-notes template create "daily" --body "# Daily Entry\n\nGrateful for:\n-\n\nAccomplishments:\n-\n\nTomorrow:\n-"

# Use template each day
apple-notes template use "daily" "Journal $(date +%Y-%m-%d)" --folder Journal

# Find today's journal entry ID
apple-notes search "Journal $(date +%Y-%m-%d)"

# Append throughout the day using the ID
apple-notes append 4318 --content "\n$(date +%H:%M): Quick note..."
```

### Meeting Notes

```bash
# Create meeting template
apple-notes template create "meeting" --body "Date: \nAttendees:\nAgenda:\n-\n\nDiscussion:\n\nAction Items:\n-"

# Create meeting note
apple-notes template use "meeting" "Weekly Team Sync" --folder Meetings

# After meeting, search for it
apple-notes search "Weekly Team Sync"
```

### Project Management

```bash
# List all project notes
apple-notes list --folder Projects

# Find notes with project tag
apple-notes tags search "#project"

# Archive completed projects
apple-notes archive --folder "Projects" --older-than 6 --to "Archive/Projects"

# Create new project note
apple-notes add "Project Alpha" --body "Goal:\n\nMilestones:\n-" --folder Projects
```

### Quick Capture

```bash
# Quick note with timestamp
apple-notes add "Quick Capture $(date +%Y-%m-%d)" --body "$(date): Idea..." --folder Inbox

# Find inbox note ID
apple-notes search "Inbox"

# Append to inbox note by ID
apple-notes append 4318 --content "\n- $(date +%H:%M): New thought"
```

## Limitations

- **macOS only**: Apple Notes database is only available on macOS
- **Plain text display**: Commands show full note content as plain text (v0.3.0+). Rich formatting (bold, italics, etc.) and media (images, attachments) are preserved in Notes.app but displayed as plain text in CLI
- **Rich content protection**: Notes with images, attachments, tables, or sketches are protected from edits by default to prevent data loss
- **Append is safer**: Use `append` instead of `edit` for notes with attachments - it won't damage existing content

## Troubleshooting

### "Database locked" error
- Close Notes.app and try again
- Notes.app may be syncing - wait a few seconds

### Note not found
- Use `apple-notes list` to see all note IDs
- Verify the note ID exists
- Use `apple-notes search "keyword"` to find the note ID

### Edit blocked (rich content detected)
- The note contains images, attachments, or other rich content
- Use Notes.app to edit, or use `apple-notes append` to add text
- Use `--force-unsafe` to override (NOT RECOMMENDED - destroys rich content)

### Notes not syncing
- Write operations use AppleScript which triggers proper sync
- Ensure you have iCloud sync enabled in Notes.app preferences
- Check internet connection for cloud sync

## Resources

- [apple-notes GitHub](https://github.com/fishfisher/apple-notes)
- [Apple Notes Database Structure](https://github.com/threeplanetssoftware/apple_cloud_notes_parser)
- Source code: `/Users/erikfisher/Developer/apple-notes`
