package database

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/isichei/recipe-book/internal/recipes"
	"github.com/isichei/recipe-book/internal/types"
)

type RecipeDatabase interface {
	SearchRecipes(string) []types.RecipeMetadata
	GetRecipeMetadata(recipeId string) types.RecipeMetadata
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

type InMemDatabase struct {
	data map[string]types.RecipeMetadata
	sc   SearchCache
}

func NewTestDatabase() InMemDatabase {
	data := map[string]types.RecipeMetadata{
		"chicken-dhansak": {
			Uid:         "chicken-dhansak",
			Title:       "Chicken Dhansak",
			Description: "A chicken dhansak recipe from BBC good foods",
		},
		"christmas-roast-potatoes": {
			Uid:         "christmas-roast-potatoes",
			Title:       "Jamie Oliver Roast Potatoes",
			Description: "A jamie oliver roast potato recipe usually used at Christmas",
		},
	}

	sc := make(SearchCache)
	sc.Add("chicken dhanksak", "chicken-dhansak")
	sc.Add("lentils", "chicken-dhansak")
	sc.Add("jamie oliver roast potatoes", "christmas-roast-potatoes")
	sc.Add("potatoes", "christmas-roast-potatoes")

	return InMemDatabase{data, sc}
}

func NewTestDatabaseFromDir(dirpath string) InMemDatabase {
	var files []string

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
	data := make(map[string]types.RecipeMetadata)
	sc := make(SearchCache)
	for _, file := range files {
		fullRecipe := recipes.ParseMarkdownFile(file)
		uid, _, _ := strings.Cut(filepath.Base(file), ".")
		uid = strings.ToLower(uid)
		d, _, _ := strings.Cut(filepath.Base(fullRecipe.Description), ".")

		// TODO: Maybe cut desc to first sentence
		data[uid] = types.RecipeMetadata{Uid: uid, Title: fullRecipe.Title, Description: d}
		for _, word := range strings.Split(strings.ToLower(fullRecipe.Title), " ") {
			sc.Add(word, uid)
		}
		for _, ingredient := range fullRecipe.Ingredients {
			sc.Add(strings.ToLower(ingredient.Name), uid)
		}
	}
	fmt.Printf("Total number of recipes read: %d\n", len(data))
	fmt.Printf("Total size of recipe cache: %d\n", len(sc))
	return InMemDatabase{data, sc}
}

// Rubbish search to fill in for a proper search query later
func (db InMemDatabase) SearchRecipes(text string) []types.RecipeMetadata {
	var filtered []types.RecipeMetadata

	if text == "" {
		max_amount := len(db.data)
		if max_amount > 9 {
			max_amount = 9
		}
		counter := 0
		for _, v := range db.data {
			filtered = append(filtered, v)
			if counter >= max_amount {
				break
			}
		}
		return filtered
	} else {
		filteredSet := make(Set)
		search_texts := strings.Split(strings.ToLower(text), " ")
		for _, search_term := range search_texts {
			for cache_key, uidSet := range db.sc {
				if strings.Contains(cache_key, search_term) {
					for uid := range uidSet {
						_, ptrs := filteredSet[uid]
						if !ptrs {
							filteredSet.Add(uid)
							filtered = append(filtered, db.data[uid])
						}
					}
				}
			}
		}
		fmt.Printf("Results count %d, filteredSet count %d\n", len(filtered), len(filteredSet))
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
