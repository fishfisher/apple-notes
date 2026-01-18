package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	// Apple Notes database path
	dbPath = "Library/Group Containers/group.com.apple.notes/NoteStore.sqlite"
)

type Note struct {
	ID        string
	Title     string
	Snippet   string
	Folder    string
	Created   time.Time
	Modified  time.Time
	Body      string
}

type Folder struct {
	Name  string
	Count int
}

type Tag struct {
	Name  string
	Count int
}

type Stats struct {
	TotalNotes      int
	TotalFolders    int
	TopTags         []Tag
	NotesThisWeek   int
	NotesThisMonth  int
	LargestNote     Note
	TotalCharacters int64
}

type DB struct {
	conn *sql.DB
}

// GetNotesDB returns the path to the Apple Notes database
func GetNotesDBPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	fullPath := filepath.Join(home, dbPath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("notes database not found at %s", fullPath)
	}

	return fullPath, nil
}

// Open opens a connection to the Apple Notes database
func Open() (*DB, error) {
	dbPath, err := GetNotesDBPath()
	if err != nil {
		return nil, err
	}

	// Open in read-only mode to avoid locking issues
	conn, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?mode=ro", dbPath))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return &DB{conn: conn}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// ListNotes retrieves all notes, optionally filtered by folder
func (db *DB) ListNotes(folder string) ([]Note, error) {
	query := `
		SELECT
			ZICCLOUDSYNCINGOBJECT.Z_PK,
			COALESCE(ZICCLOUDSYNCINGOBJECT.ZTITLE1, '') as title,
			COALESCE(ZICCLOUDSYNCINGOBJECT.ZSNIPPET, '') as snippet,
			COALESCE(folders.ZTITLE2, 'Notes') as folder,
			COALESCE(datetime(ZICCLOUDSYNCINGOBJECT.ZCREATIONDATE + 978307200, 'unixepoch', 'localtime'), '') as created,
			COALESCE(datetime(ZICCLOUDSYNCINGOBJECT.ZMODIFICATIONDATE1 + 978307200, 'unixepoch', 'localtime'), '') as modified
		FROM ZICCLOUDSYNCINGOBJECT
		LEFT JOIN ZICCLOUDSYNCINGOBJECT as folders ON ZICCLOUDSYNCINGOBJECT.ZFOLDER = folders.Z_PK
		WHERE ZICCLOUDSYNCINGOBJECT.ZTITLE1 IS NOT NULL
			AND ZICCLOUDSYNCINGOBJECT.ZMARKEDFORDELETION = 0
	`

	var args []interface{}
	if folder != "" {
		query += " AND folders.ZTITLE2 = ?"
		args = append(args, folder)
	}

	query += " ORDER BY ZICCLOUDSYNCINGOBJECT.ZMODIFICATIONDATE1 DESC"

	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query notes: %w", err)
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		var createdStr, modifiedStr string
		err := rows.Scan(&note.ID, &note.Title, &note.Snippet, &note.Folder, &createdStr, &modifiedStr)
		if err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}

		if createdStr != "" {
			note.Created, _ = time.Parse("2006-01-02 15:04:05", createdStr)
		}
		if modifiedStr != "" {
			note.Modified, _ = time.Parse("2006-01-02 15:04:05", modifiedStr)
		}
		notes = append(notes, note)
	}

	return notes, nil
}

// SearchNotes searches for notes containing the search term in title or snippet
func (db *DB) SearchNotes(term string) ([]Note, error) {
	query := `
		SELECT
			ZICCLOUDSYNCINGOBJECT.Z_PK,
			COALESCE(ZICCLOUDSYNCINGOBJECT.ZTITLE1, '') as title,
			COALESCE(ZICCLOUDSYNCINGOBJECT.ZSNIPPET, '') as snippet,
			COALESCE(folders.ZTITLE2, 'Notes') as folder,
			COALESCE(datetime(ZICCLOUDSYNCINGOBJECT.ZCREATIONDATE + 978307200, 'unixepoch', 'localtime'), '') as created,
			COALESCE(datetime(ZICCLOUDSYNCINGOBJECT.ZMODIFICATIONDATE1 + 978307200, 'unixepoch', 'localtime'), '') as modified
		FROM ZICCLOUDSYNCINGOBJECT
		LEFT JOIN ZICCLOUDSYNCINGOBJECT as folders ON ZICCLOUDSYNCINGOBJECT.ZFOLDER = folders.Z_PK
		WHERE ZICCLOUDSYNCINGOBJECT.ZTITLE1 IS NOT NULL
			AND ZICCLOUDSYNCINGOBJECT.ZMARKEDFORDELETION = 0
			AND (ZICCLOUDSYNCINGOBJECT.ZTITLE1 LIKE ? OR ZICCLOUDSYNCINGOBJECT.ZSNIPPET LIKE ?)
		ORDER BY ZICCLOUDSYNCINGOBJECT.ZMODIFICATIONDATE1 DESC
		LIMIT 100
	`

	searchPattern := "%" + term + "%"
	rows, err := db.conn.Query(query, searchPattern, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search notes: %w", err)
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		var createdStr, modifiedStr string
		err := rows.Scan(&note.ID, &note.Title, &note.Snippet, &note.Folder, &createdStr, &modifiedStr)
		if err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}

		if createdStr != "" {
			note.Created, _ = time.Parse("2006-01-02 15:04:05", createdStr)
		}
		if modifiedStr != "" {
			note.Modified, _ = time.Parse("2006-01-02 15:04:05", modifiedStr)
		}
		notes = append(notes, note)
	}

	return notes, nil
}

// GetNote retrieves a specific note by ID or title
func (db *DB) GetNote(identifier string) (*Note, error) {
	query := `
		SELECT
			ZICCLOUDSYNCINGOBJECT.Z_PK,
			COALESCE(ZICCLOUDSYNCINGOBJECT.ZTITLE1, '') as title,
			COALESCE(ZICCLOUDSYNCINGOBJECT.ZSNIPPET, '') as snippet,
			COALESCE(folders.ZTITLE2, 'Notes') as folder,
			COALESCE(datetime(ZICCLOUDSYNCINGOBJECT.ZCREATIONDATE + 978307200, 'unixepoch', 'localtime'), '') as created,
			COALESCE(datetime(ZICCLOUDSYNCINGOBJECT.ZMODIFICATIONDATE1 + 978307200, 'unixepoch', 'localtime'), '') as modified
		FROM ZICCLOUDSYNCINGOBJECT
		LEFT JOIN ZICCLOUDSYNCINGOBJECT as folders ON ZICCLOUDSYNCINGOBJECT.ZFOLDER = folders.Z_PK
		WHERE ZICCLOUDSYNCINGOBJECT.ZMARKEDFORDELETION = 0
			AND (ZICCLOUDSYNCINGOBJECT.Z_PK = ? OR ZICCLOUDSYNCINGOBJECT.ZTITLE1 = ?)
		LIMIT 1
	`

	var note Note
	var createdStr, modifiedStr string
	err := db.conn.QueryRow(query, identifier, identifier).Scan(
		&note.ID, &note.Title, &note.Snippet, &note.Folder, &createdStr, &modifiedStr,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("note not found: %s", identifier)
		}
		return nil, fmt.Errorf("failed to get note: %w", err)
	}

	if createdStr != "" {
		note.Created, _ = time.Parse("2006-01-02 15:04:05", createdStr)
	}
	if modifiedStr != "" {
		note.Modified, _ = time.Parse("2006-01-02 15:04:05", modifiedStr)
	}

	return &note, nil
}

// ListFolders retrieves all note folders with note counts
func (db *DB) ListFolders() ([]Folder, error) {
	query := `
		SELECT
			COALESCE(folders.ZTITLE2, 'Notes') as folder_name,
			COUNT(notes.Z_PK) as note_count
		FROM ZICCLOUDSYNCINGOBJECT as notes
		LEFT JOIN ZICCLOUDSYNCINGOBJECT as folders ON notes.ZFOLDER = folders.Z_PK
		WHERE notes.ZTITLE1 IS NOT NULL
			AND notes.ZMARKEDFORDELETION = 0
		GROUP BY folder_name
		ORDER BY folder_name
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query folders: %w", err)
	}
	defer rows.Close()

	var folders []Folder
	for rows.Next() {
		var folder Folder
		err := rows.Scan(&folder.Name, &folder.Count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan folder: %w", err)
		}
		folders = append(folders, folder)
	}

	return folders, nil
}

// GetRecentNotes retrieves notes modified within the specified number of days
func (db *DB) GetRecentNotes(days int, limit int) ([]Note, error) {
	query := `
		SELECT
			ZICCLOUDSYNCINGOBJECT.Z_PK,
			COALESCE(ZICCLOUDSYNCINGOBJECT.ZTITLE1, '') as title,
			COALESCE(ZICCLOUDSYNCINGOBJECT.ZSNIPPET, '') as snippet,
			COALESCE(folders.ZTITLE2, 'Notes') as folder,
			COALESCE(datetime(ZICCLOUDSYNCINGOBJECT.ZCREATIONDATE + 978307200, 'unixepoch', 'localtime'), '') as created,
			COALESCE(datetime(ZICCLOUDSYNCINGOBJECT.ZMODIFICATIONDATE1 + 978307200, 'unixepoch', 'localtime'), '') as modified
		FROM ZICCLOUDSYNCINGOBJECT
		LEFT JOIN ZICCLOUDSYNCINGOBJECT as folders ON ZICCLOUDSYNCINGOBJECT.ZFOLDER = folders.Z_PK
		WHERE ZICCLOUDSYNCINGOBJECT.ZTITLE1 IS NOT NULL
			AND ZICCLOUDSYNCINGOBJECT.ZMARKEDFORDELETION = 0
			AND ZICCLOUDSYNCINGOBJECT.ZMODIFICATIONDATE1 > (strftime('%s', 'now') - 978307200 - (? * 86400))
		ORDER BY ZICCLOUDSYNCINGOBJECT.ZMODIFICATIONDATE1 DESC
		LIMIT ?
	`

	rows, err := db.conn.Query(query, days, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent notes: %w", err)
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		var createdStr, modifiedStr string
		err := rows.Scan(&note.ID, &note.Title, &note.Snippet, &note.Folder, &createdStr, &modifiedStr)
		if err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}

		if createdStr != "" {
			note.Created, _ = time.Parse("2006-01-02 15:04:05", createdStr)
		}
		if modifiedStr != "" {
			note.Modified, _ = time.Parse("2006-01-02 15:04:05", modifiedStr)
		}
		notes = append(notes, note)
	}

	return notes, nil
}

// ExtractTags extracts all hashtags from note snippets
func (db *DB) ExtractTags() ([]Tag, error) {
	query := `
		SELECT ZSNIPPET
		FROM ZICCLOUDSYNCINGOBJECT
		WHERE ZTITLE1 IS NOT NULL
			AND ZMARKEDFORDELETION = 0
			AND ZSNIPPET IS NOT NULL
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query notes for tags: %w", err)
	}
	defer rows.Close()

	tagCounts := make(map[string]int)
	for rows.Next() {
		var snippet string
		if err := rows.Scan(&snippet); err != nil {
			continue
		}

		// Extract hashtags from snippet
		tags := extractHashtags(snippet)
		for _, tag := range tags {
			tagCounts[tag]++
		}
	}

	// Convert map to slice
	var tags []Tag
	for name, count := range tagCounts {
		tags = append(tags, Tag{Name: name, Count: count})
	}

	// Sort by count descending
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Count > tags[j].Count
	})

	return tags, nil
}

// SearchByTag searches for notes containing a specific hashtag
func (db *DB) SearchByTag(tag string) ([]Note, error) {
	// Ensure tag starts with #
	if !strings.HasPrefix(tag, "#") {
		tag = "#" + tag
	}

	return db.SearchNotes(tag)
}

// GetStats returns statistics about the notes collection
func (db *DB) GetStats() (*Stats, error) {
	stats := &Stats{}

	// Total notes
	err := db.conn.QueryRow(`
		SELECT COUNT(*)
		FROM ZICCLOUDSYNCINGOBJECT
		WHERE ZTITLE1 IS NOT NULL AND ZMARKEDFORDELETION = 0
	`).Scan(&stats.TotalNotes)
	if err != nil {
		return nil, fmt.Errorf("failed to count notes: %w", err)
	}

	// Total folders
	err = db.conn.QueryRow(`
		SELECT COUNT(DISTINCT ZFOLDER)
		FROM ZICCLOUDSYNCINGOBJECT
		WHERE ZTITLE1 IS NOT NULL AND ZMARKEDFORDELETION = 0
	`).Scan(&stats.TotalFolders)
	if err != nil {
		return nil, fmt.Errorf("failed to count folders: %w", err)
	}

	// Notes this week (7 days)
	err = db.conn.QueryRow(`
		SELECT COUNT(*)
		FROM ZICCLOUDSYNCINGOBJECT
		WHERE ZTITLE1 IS NOT NULL
			AND ZMARKEDFORDELETION = 0
			AND ZCREATIONDATE > (strftime('%s', 'now') - 978307200 - (7 * 86400))
	`).Scan(&stats.NotesThisWeek)
	if err != nil {
		return nil, fmt.Errorf("failed to count weekly notes: %w", err)
	}

	// Notes this month (30 days)
	err = db.conn.QueryRow(`
		SELECT COUNT(*)
		FROM ZICCLOUDSYNCINGOBJECT
		WHERE ZTITLE1 IS NOT NULL
			AND ZMARKEDFORDELETION = 0
			AND ZCREATIONDATE > (strftime('%s', 'now') - 978307200 - (30 * 86400))
	`).Scan(&stats.NotesThisMonth)
	if err != nil {
		return nil, fmt.Errorf("failed to count monthly notes: %w", err)
	}

	// Largest note
	var title, snippet, folder string
	var createdStr, modifiedStr string
	err = db.conn.QueryRow(`
		SELECT
			COALESCE(ZTITLE1, '') as title,
			COALESCE(ZSNIPPET, '') as snippet,
			COALESCE((SELECT ZTITLE2 FROM ZICCLOUDSYNCINGOBJECT WHERE Z_PK = n.ZFOLDER), 'Notes') as folder,
			COALESCE(datetime(ZCREATIONDATE + 978307200, 'unixepoch', 'localtime'), '') as created,
			COALESCE(datetime(ZMODIFICATIONDATE1 + 978307200, 'unixepoch', 'localtime'), '') as modified,
			LENGTH(ZSNIPPET) as len
		FROM ZICCLOUDSYNCINGOBJECT n
		WHERE ZTITLE1 IS NOT NULL
			AND ZMARKEDFORDELETION = 0
			AND ZSNIPPET IS NOT NULL
		ORDER BY len DESC
		LIMIT 1
	`).Scan(&title, &snippet, &folder, &createdStr, &modifiedStr, &stats.TotalCharacters)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to find largest note: %w", err)
	}
	if err == nil {
		stats.LargestNote = Note{Title: title, Snippet: snippet, Folder: folder}
		if createdStr != "" {
			stats.LargestNote.Created, _ = time.Parse("2006-01-02 15:04:05", createdStr)
		}
		if modifiedStr != "" {
			stats.LargestNote.Modified, _ = time.Parse("2006-01-02 15:04:05", modifiedStr)
		}
	}

	// Top tags
	tags, err := db.ExtractTags()
	if err == nil && len(tags) > 10 {
		stats.TopTags = tags[:10]
	} else {
		stats.TopTags = tags
	}

	return stats, nil
}

// FindDuplicates finds notes with identical or very similar titles
func (db *DB) FindDuplicates() ([][]Note, error) {
	query := `
		SELECT
			ZICCLOUDSYNCINGOBJECT.Z_PK,
			COALESCE(ZICCLOUDSYNCINGOBJECT.ZTITLE1, '') as title,
			COALESCE(ZICCLOUDSYNCINGOBJECT.ZSNIPPET, '') as snippet,
			COALESCE(folders.ZTITLE2, 'Notes') as folder,
			COALESCE(datetime(ZICCLOUDSYNCINGOBJECT.ZCREATIONDATE + 978307200, 'unixepoch', 'localtime'), '') as created,
			COALESCE(datetime(ZICCLOUDSYNCINGOBJECT.ZMODIFICATIONDATE1 + 978307200, 'unixepoch', 'localtime'), '') as modified
		FROM ZICCLOUDSYNCINGOBJECT
		LEFT JOIN ZICCLOUDSYNCINGOBJECT as folders ON ZICCLOUDSYNCINGOBJECT.ZFOLDER = folders.Z_PK
		WHERE ZICCLOUDSYNCINGOBJECT.ZTITLE1 IS NOT NULL
			AND ZICCLOUDSYNCINGOBJECT.ZMARKEDFORDELETION = 0
		ORDER BY ZICCLOUDSYNCINGOBJECT.ZTITLE1
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query notes: %w", err)
	}
	defer rows.Close()

	titleMap := make(map[string][]Note)
	for rows.Next() {
		var note Note
		var createdStr, modifiedStr string
		err := rows.Scan(&note.ID, &note.Title, &note.Snippet, &note.Folder, &createdStr, &modifiedStr)
		if err != nil {
			continue
		}

		if createdStr != "" {
			note.Created, _ = time.Parse("2006-01-02 15:04:05", createdStr)
		}
		if modifiedStr != "" {
			note.Modified, _ = time.Parse("2006-01-02 15:04:05", modifiedStr)
		}

		titleMap[note.Title] = append(titleMap[note.Title], note)
	}

	// Find duplicates
	var duplicates [][]Note
	for _, notes := range titleMap {
		if len(notes) > 1 {
			duplicates = append(duplicates, notes)
		}
	}

	return duplicates, nil
}

// ExtractLinks extracts all URLs from a note's snippet
func (db *DB) ExtractLinks(noteIdentifier string) ([]string, error) {
	note, err := db.GetNote(noteIdentifier)
	if err != nil {
		return nil, err
	}

	return extractURLs(note.Snippet), nil
}

// FindNotesWithLinks finds all notes that contain URLs
func (db *DB) FindNotesWithLinks() ([]Note, error) {
	query := `
		SELECT
			ZICCLOUDSYNCINGOBJECT.Z_PK,
			COALESCE(ZICCLOUDSYNCINGOBJECT.ZTITLE1, '') as title,
			COALESCE(ZICCLOUDSYNCINGOBJECT.ZSNIPPET, '') as snippet,
			COALESCE(folders.ZTITLE2, 'Notes') as folder,
			COALESCE(datetime(ZICCLOUDSYNCINGOBJECT.ZCREATIONDATE + 978307200, 'unixepoch', 'localtime'), '') as created,
			COALESCE(datetime(ZICCLOUDSYNCINGOBJECT.ZMODIFICATIONDATE1 + 978307200, 'unixepoch', 'localtime'), '') as modified
		FROM ZICCLOUDSYNCINGOBJECT
		LEFT JOIN ZICCLOUDSYNCINGOBJECT as folders ON ZICCLOUDSYNCINGOBJECT.ZFOLDER = folders.Z_PK
		WHERE ZICCLOUDSYNCINGOBJECT.ZTITLE1 IS NOT NULL
			AND ZICCLOUDSYNCINGOBJECT.ZMARKEDFORDELETION = 0
			AND (ZICCLOUDSYNCINGOBJECT.ZSNIPPET LIKE '%http://%' OR ZICCLOUDSYNCINGOBJECT.ZSNIPPET LIKE '%https://%')
		ORDER BY ZICCLOUDSYNCINGOBJECT.ZMODIFICATIONDATE1 DESC
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query notes with links: %w", err)
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		var createdStr, modifiedStr string
		err := rows.Scan(&note.ID, &note.Title, &note.Snippet, &note.Folder, &createdStr, &modifiedStr)
		if err != nil {
			continue
		}

		if createdStr != "" {
			note.Created, _ = time.Parse("2006-01-02 15:04:05", createdStr)
		}
		if modifiedStr != "" {
			note.Modified, _ = time.Parse("2006-01-02 15:04:05", modifiedStr)
		}

		notes = append(notes, note)
	}

	return notes, nil
}

// Helper functions

func extractHashtags(text string) []string {
	var tags []string
	words := strings.Fields(text)
	for _, word := range words {
		if strings.HasPrefix(word, "#") {
			// Clean up the tag (remove punctuation at the end)
			tag := strings.TrimRight(word, ".,!?;:)")
			if len(tag) > 1 {
				tags = append(tags, tag)
			}
		}
	}
	return tags
}

func extractURLs(text string) []string {
	var urls []string
	words := strings.Fields(text)
	for _, word := range words {
		if strings.HasPrefix(word, "http://") || strings.HasPrefix(word, "https://") {
			urls = append(urls, word)
		}
	}
	return urls
}
