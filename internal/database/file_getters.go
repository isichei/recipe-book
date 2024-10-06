package database

import (
	"fmt"
	"github.com/isichei/recipe-book/internal/recipes"
	"os"
	"path/filepath"
	"strings"
)

type MarkdownFileGetter interface {
	getRecipe(string) recipes.Recipe
	files() []string
	getRecipeFromFilePath(string) (string, recipes.Recipe)
}

// LOCAL MARKDOWN GETTER
type LocalMarkdownFileGetter struct {
	dir string
}

func (fg LocalMarkdownFileGetter) files() []string {
	var files []string

	err := filepath.Walk(fg.dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") && !strings.HasSuffix(info.Name(), "template.md") {
			fmt.Println(path)
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		panic("Failed to read files in dir cannot create db")
	}

	if len(files) == 0 {
		panic("No files found")
	}
	return files
}

// Get a recipe from the local file system
func (fg LocalMarkdownFileGetter) getRecipe(uid string) recipes.Recipe {
	data := read_file(fg.dir + uid + ".md")
	return recipes.ParseMarkdownFile(data)
}

// Get a recipe from the local file system based on the file path, return uid and recipe
func (fg LocalMarkdownFileGetter) getRecipeFromFilePath(file string) (string, recipes.Recipe) {
	uid := getUidFromFilePath(file)
	return uid, fg.getRecipe(uid)
}

func getUidFromFilePath(file string) string {
	uid, _, _ := strings.Cut(filepath.Base(file), ".")
	return uid
}

func read_file(filepath string) string {
	data, err := os.ReadFile(filepath)
	if err != nil {
		msg := fmt.Sprintf("Error opening file: %s. %s", filepath, err)
		fmt.Println(msg)
		panic(msg)
	}

	// Create a scanner to read the file line by line
	return string(data)
}
