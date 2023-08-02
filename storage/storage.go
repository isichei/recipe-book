package storage

import (
	"strings"

	"github.com/isichei/recipe-book/types"
)

type Storage interface {
	SearchRecipes(string) []types.RecipeMetadata
}

type FakeStorage struct {
	data []types.RecipeMetadata
}

func NewFakeStorage() FakeStorage {
	dhansak := types.RecipeMetadata{
		Uid:         "chicken-dhansak-recipe",
		Title:       "Chicken Dhansak",
		Description: "A chicken dhansak recipe from BBC good foods",
	}
	roast := types.RecipeMetadata{
		Uid:         "christmas-roast-potatoes",
		Title:       "Jamie Oliver Roast Potatoes",
		Description: "A jamie oliver roast potato recipe usually used at Christmas",
	}
	return FakeStorage{
		data: []types.RecipeMetadata{dhansak, roast},
	}
}

func (store FakeStorage) SearchRecipes(text string) []types.RecipeMetadata {

	if text == "" {
		return store.data
	} else {
		var filtered []types.RecipeMetadata

		for _, recipe := range store.data {
			if strings.Contains(strings.ToLower(recipe.Description), strings.ToLower(text)) {
				filtered = append(filtered, recipe)
			}
		}
		return filtered
	}
}
