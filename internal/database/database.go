package database

import (
	"github.com/isichei/recipe-book/internal/types"
	"strings"
)

type RecipeDatabase interface {
	SearchRecipes(string) []types.RecipeMetadata
}

type TestDatabase struct {
	data []types.RecipeMetadata
}

func NewTestDatabase() TestDatabase {
	data := []types.RecipeMetadata{
		{
			Uid:         "chicken-dhansak-recipe",
			Title:       "Chicken Dhansak",
			Description: "A chicken dhansak recipe from BBC good foods",
		},
		{
			Uid:         "christmas-roast-potatoes",
			Title:       "Jamie Oliver Roast Potatoes",
			Description: "A jamie oliver roast potato recipe usually used at Christmas",
		},
	}
	return TestDatabase{data}
}

func filterReceipeMetadata(recipes []types.RecipeMetadata, text string) []types.RecipeMetadata {
	var filtered []types.RecipeMetadata

	for _, recipe := range recipes {
		if strings.Contains(strings.ToLower(recipe.Description), strings.ToLower(text)) {
			filtered = append(filtered, recipe)
		}
	}

	return filtered
}

// Rubbish search to fill in for a proper search query later
func (db TestDatabase) SearchRecipes(text string) []types.RecipeMetadata {
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
