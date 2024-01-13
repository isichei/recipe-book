package database

import (
	"github.com/isichei/recipe-book/internal/recipes"
	"github.com/isichei/recipe-book/internal/types"
	"os"
	"path/filepath"
	"strings"
)

type RecipeDatabase interface {
	SearchRecipes(string) []types.RecipeMetadata
	GetRecipeMetadata(recipeId string) types.RecipeMetadata
}

type InMemDatabase struct {
	data []types.RecipeMetadata
}

func NewTestDatabase() InMemDatabase {
	data := []types.RecipeMetadata{
		{
			Uid:         "chicken-dhansak",
			Title:       "Chicken Dhansak",
			Description: "A chicken dhansak recipe from BBC good foods",
		},
		{
			Uid:         "christmas-roast-potatoes",
			Title:       "Jamie Oliver Roast Potatoes",
			Description: "A jamie oliver roast potato recipe usually used at Christmas",
		},
	}
	return InMemDatabase{data}
}

func NewTestDatabaseFromDir(dirpath string) InMemDatabase {
	var files []string
	var data []types.RecipeMetadata

	err := filepath.Walk(dirpath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") && !strings.HasSuffix(info.Name(), "template.md") {
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

	for _, file := range files {
		fullRecipe := recipes.ParseMarkdownFile(file)
		uid, _, _ := strings.Cut(filepath.Base(file), ".")
		d, _, _ := strings.Cut(filepath.Base(fullRecipe.Description), ".")
		// TODO: Maybe cut desc to first sentence
		data = append(data, types.RecipeMetadata{Uid: uid, Title: fullRecipe.Title, Description: d})
	}
	return InMemDatabase{data}
}

// Rubbish search to fill in for a proper search query later
func (db InMemDatabase) SearchRecipes(text string) []types.RecipeMetadata {
	if text == "" {
		return db.data
	} else {
		var filtered []types.RecipeMetadata

		for _, recipe := range db.data {
			if strings.Contains(strings.ToLower(recipe.Description), strings.ToLower(text)) {
				filtered = append(filtered, recipe)
			}
		}
		return filtered
	}
}

func (db InMemDatabase) GetRecipeMetadata(recipeUid string) types.RecipeMetadata {
	for _, d := range db.data {
		if d.Uid == recipeUid {
			return d
		}
	}
	return types.RecipeMetadata{}
}
