package filesyncer

import (
	"errors"
	"fmt"
	"strings"

	"github.com/isichei/recipe-book/internal/database"
	"github.com/isichei/recipe-book/internal/recipes"
	"log/slog"
	"os"
	"path"
)

// Light interface to interact with files or db
// string -> filename , []byte -> data
type RecipeReadWriter interface {
	Write(string, []byte) error  //
	Read(string) ([]byte, error) // string -> filename
	Delete(string) error
}

type RawMdRecipeReadWriter struct {
	directory string
}

func NewRawMdRecipeReadWriter(directory string) *RawMdRecipeReadWriter {
	return &RawMdRecipeReadWriter{directory: directory}
}

func (rrw *RawMdRecipeReadWriter) Write(filename string, data []byte) error {
	err := os.WriteFile(path.Join(rrw.directory, filename), data, 0644)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to write %s from msg", filename), err)
	}
	return nil
}

func (rrw *RawMdRecipeReadWriter) Read(filename string) ([]byte, error) {
	return os.ReadFile(path.Join(rrw.directory, filename))
}

func (rrw *RawMdRecipeReadWriter) Delete(filename string) error {
	fileToDelete := path.Join(rrw.directory, filename)
	err := os.Remove(fileToDelete)
	if err != nil {
		slog.Error("Replica could not delete recipe", "filename", filename, "error", err)
		return fmt.Errorf("Replica failed to delete file %s: %w", filename, err)
	} else {
		slog.Debug("Replica deleting", "filename", filename)
	}
	return nil
}

type DbRecipeReadWriter struct {
	db database.RecipeDatabase
}

func (rrw *DbRecipeReadWriter) Write(filename string, data []byte) error {
	rUid, err := filenameToUid(filename)
	if err != nil {
		return err
	}
	recipe := recipes.ParseMarkdownFile(string(data))
	return rrw.db.AddRecipe(rUid, recipe)
}

func (rrw *DbRecipeReadWriter) Read(filename string) ([]byte, error) {
	panic("Not implemented as a file reader")
}

func (rrw *DbRecipeReadWriter) Delete(filename string) error {
	rUid, err := filenameToUid(filename)
	if err != nil {
		return err
	}
	return rrw.db.DeleteRecipe(rUid)
}

func filenameToUid(filename string) (string, error) {
	rUid, _, ok := strings.Cut(filename, ".")
	if !ok {
		return "", fmt.Errorf("Implementation error, expecting msg.Filename to be something like: `baked-beans.md` but got %s from msg", filename)
	}
	return rUid, nil
}
