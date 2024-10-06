package database

import (
	"github.com/isichei/recipe-book/internal/recipes"
)

type RecipeDatabase interface {
	SearchRecipes(string) []recipes.RecipeMetadata
	GetRecipeMetadata(recipeId string) recipes.RecipeMetadata
	GetRecipe(recipeUid string) recipes.Recipe
	AddRecipe(rUid string, r recipes.Recipe) error
}
type Set map[string]bool

func (s Set) Add(item string) {
	s[item] = true
}

type SearchCache map[string]Set

func (sc SearchCache) Add(item string, uid string) {
	set, ptrs := sc[item]
	if !ptrs {
		set = make(Set)
	}
	set.Add(uid)
	sc[item] = set
}

func (sc SearchCache) Retrieve(item string) Set {
	return sc[item]
}
