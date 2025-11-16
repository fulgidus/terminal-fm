// Package storage provides database functionality for bookmarks and user data.
package storage

import (
	"database/sql"
	"fmt"

	"github.com/fulgidus/terminal-fm/pkg/services/radiobrowser"
	_ "github.com/mattn/go-sqlite3"
)

// Store handles database operations.
type Store struct {
	db *sql.DB
}

// NewStore creates a new database store and initializes the schema.
func NewStore(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	store := &Store{db: db}

	// Initialize schema
	if err := store.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return store, nil
}

// Close closes the database connection.
func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// initSchema creates the necessary tables.
func (s *Store) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS bookmarks (
		station_uuid TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		url TEXT NOT NULL,
		url_resolved TEXT NOT NULL,
		homepage TEXT,
		tags TEXT,
		country TEXT,
		country_code TEXT,
		language TEXT,
		language_codes TEXT,
		votes INTEGER,
		codec TEXT,
		bitrate INTEGER,
		last_check_ok INTEGER,
		click_count INTEGER,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE INDEX IF NOT EXISTS idx_bookmarks_name ON bookmarks(name);
	CREATE INDEX IF NOT EXISTS idx_bookmarks_created ON bookmarks(created_at);
	`

	_, err := s.db.Exec(schema)
	return err
}

// AddBookmark adds a station to bookmarks.
func (s *Store) AddBookmark(station *radiobrowser.Station) error {
	if station == nil {
		return fmt.Errorf("station cannot be nil")
	}

	query := `
	INSERT INTO bookmarks (
		station_uuid, name, url, url_resolved, homepage, tags,
		country, country_code, language, language_codes,
		votes, codec, bitrate, last_check_ok, click_count
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(station_uuid) DO NOTHING
	`

	_, err := s.db.Exec(query,
		station.StationUUID,
		station.Name,
		station.URL,
		station.URLResolved,
		station.Homepage,
		station.Tags,
		station.Country,
		station.CountryCode,
		station.Language,
		station.LanguageCodes,
		station.Votes,
		station.Codec,
		station.Bitrate,
		station.LastCheckOK,
		station.ClickCount,
	)

	if err != nil {
		return fmt.Errorf("failed to add bookmark: %w", err)
	}

	return nil
}

// RemoveBookmark removes a station from bookmarks.
func (s *Store) RemoveBookmark(stationUUID string) error {
	query := `DELETE FROM bookmarks WHERE station_uuid = ?`

	result, err := s.db.Exec(query, stationUUID)
	if err != nil {
		return fmt.Errorf("failed to remove bookmark: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("bookmark not found")
	}

	return nil
}

// GetBookmarks retrieves all bookmarked stations.
func (s *Store) GetBookmarks() ([]radiobrowser.Station, error) {
	query := `
	SELECT 
		station_uuid, name, url, url_resolved, homepage, tags,
		country, country_code, language, language_codes,
		votes, codec, bitrate, last_check_ok, click_count
	FROM bookmarks
	ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query bookmarks: %w", err)
	}
	defer rows.Close()

	var bookmarks []radiobrowser.Station

	for rows.Next() {
		var station radiobrowser.Station
		err := rows.Scan(
			&station.StationUUID,
			&station.Name,
			&station.URL,
			&station.URLResolved,
			&station.Homepage,
			&station.Tags,
			&station.Country,
			&station.CountryCode,
			&station.Language,
			&station.LanguageCodes,
			&station.Votes,
			&station.Codec,
			&station.Bitrate,
			&station.LastCheckOK,
			&station.ClickCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bookmark: %w", err)
		}

		bookmarks = append(bookmarks, station)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating bookmarks: %w", err)
	}

	return bookmarks, nil
}

// IsBookmarked checks if a station is bookmarked.
func (s *Store) IsBookmarked(stationUUID string) (bool, error) {
	query := `SELECT COUNT(*) FROM bookmarks WHERE station_uuid = ?`

	var count int
	err := s.db.QueryRow(query, stationUUID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check bookmark: %w", err)
	}

	return count > 0, nil
}

// GetBookmarkCount returns the total number of bookmarks.
func (s *Store) GetBookmarkCount() (int, error) {
	query := `SELECT COUNT(*) FROM bookmarks`

	var count int
	err := s.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count bookmarks: %w", err)
	}

	return count, nil
}

// CleanOldBackups removes bookmark backups older than specified days.
func (s *Store) CleanOldBackups(days int) error {
	// TODO: Implement backup cleanup in future version
	return nil
}

// ExportBookmarks exports bookmarks to JSON (for future backup feature).
func (s *Store) ExportBookmarks() ([]byte, error) {
	// TODO: Implement JSON export in future version
	return nil, fmt.Errorf("not implemented")
}

// ImportBookmarks imports bookmarks from JSON (for future backup feature).
func (s *Store) ImportBookmarks(data []byte) error {
	// TODO: Implement JSON import in future version
	return fmt.Errorf("not implemented")
}
