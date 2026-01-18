package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
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
