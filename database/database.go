package database

import "strings"

type RecipeMetadata struct {
	Uid         string
	Title       string
	Description string
}

type RecipeDatabase interface {
	SearchRecipes(string) []RecipeMetadata
}

type TestDatabase struct{}

func filterReceipeMetadata(recipes []RecipeMetadata, text string) []RecipeMetadata {
	var filtered []RecipeMetadata

	for _, recipe := range recipes {
		if strings.Contains(strings.ToLower(recipe.Description), strings.ToLower(text)) {
			filtered = append(filtered, recipe)
		}
	}

	return filtered
}

// Rubbish search to fill in for a proper search query later
func (db *TestDatabase) SearchRecipes(text string) []RecipeMetadata {

	data := []RecipeMetadata{
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

	if text == "" {
		return data
	} else {
		var filtered []RecipeMetadata

		for _, recipe := range data {
			if strings.Contains(strings.ToLower(recipe.Description), strings.ToLower(text)) {
				filtered = append(filtered, recipe)
			}
		}
		return filtered
	}
}
