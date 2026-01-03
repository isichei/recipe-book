package filesyncer

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"iter"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type fileCacheData struct {
	md5    string
	synced bool
}

type FileCache interface {
	All() iter.Seq2[string, fileCacheData] // string is filename
	GetDirectory() string                  // TODO: See if can get rid
	Get(string) (fileCacheData, bool)      // bool is exists
	Add(string, fileCacheData)             // string is filename
}

type RawMdFileCache struct {
	data      map[string]fileCacheData
	directory string
}

func (fc *RawMdFileCache) All() iter.Seq2[string, fileCacheData] {
	return func(yield func(string, fileCacheData) bool) {
		for k, v := range fc.data {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (fc *RawMdFileCache) GetDirectory() string {
	return fc.directory
}

func (fc *RawMdFileCache) Get(filename string) (fileCacheData, bool) {
	fcd, ok := fc.data[filename]
	return fcd, ok
}

func (fc *RawMdFileCache) Add(filename string, fileData fileCacheData) {
	fc.data[filename] = fileData
}

type DbFileCache struct {
	dbEngine *sql.DB
}

func (fc *DbFileCache) All() iter.Seq2[string, fileCacheData] {
	return func(yield func(string, fileCacheData) bool) {
		query := `SELECT id, md5, synced
		FROM file_cache
		WHERE deleted = 0;
		`
		rows, err := fc.dbEngine.Query(query)
		if err != nil {
			panic("Could no query database when trying to get file_cache") // TODO more thought needed here
		}
		defer rows.Close()
		for rows.Next() {
			var recipe_id string
			var fcd fileCacheData

			rows.Scan(&recipe_id, &fcd.md5, &fcd.synced)
			if !yield(recipe_id, fcd) {
				return
			}
		}
	}
}

func (fc *DbFileCache) GetDirectory() string {
	panic("Not implemented")
}

func (fc *DbFileCache) Get(filename string) (fileCacheData, bool) {
	query := `SELECT md5, synced
		FROM file_cache
		WHERE id = ? and deleted = 0;`

	var fcd fileCacheData
	err := fc.dbEngine.QueryRow(query, filename).Scan(&fcd.md5, &fcd.synced)
	if err != nil {
		if err == sql.ErrNoRows {
			return fcd, false
		} else {
			panic(fmt.Sprintf("DB error on Get fileCacheData for %s", filename))
		}
	}
	return fcd, true
}

func (fc *DbFileCache) Add(filename string, fileData fileCacheData) {
	query := `INSERT INTO file_cache (id, md5, deleted, last_edited, synced)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		md5 = excluded.md5,
		deleted = excluded.deleted,
		last_edited = excluded.last_edited,
		synced = excluded.synced;
	`
	current_time := time.Now().UTC().Format(time.DateTime)
	_, err := fc.dbEngine.Exec(
		query,
		filename,
		fileData.md5,
		false,
		current_time,
		fileData.synced,
	)
	if err != nil {
		panic(fmt.Sprintf("Could not add fileCacheData to the DB - %s", err))
	}
}

// Returns a RawMdFileCache with md files from a given directory scanned
func CreateRawMdFileCache(directory string) (*RawMdFileCache, error) {
	fc := RawMdFileCache{directory: directory, data: map[string]fileCacheData{}}

	c, err := os.ReadDir(directory)
	if err != nil {
		return nil, errors.Join(errors.New("Failed to open directory"), err)
	}

	for _, entry := range c {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		f, err := os.Open(filepath.Join(directory, entry.Name()))
		if err != nil {
			slog.Error("Failed to open file", "error", err)
			return nil, fmt.Errorf("Failed to open file %s: %w", entry.Name(), err)
		}

		h := md5.New()
		if _, err := io.Copy(h, f); err != nil {
			slog.Error("Failed to hash file", "error", err)
			f.Close()
			return nil, fmt.Errorf("Failed to hash file %s: %w", entry.Name(), err)
		}

		hash := hex.EncodeToString(h.Sum(nil))
		fc.data[entry.Name()] = fileCacheData{md5: hash, synced: false}
		f.Close()
	}
	return &fc, nil
}
