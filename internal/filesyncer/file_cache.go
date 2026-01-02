package filesyncer

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"iter"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
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
