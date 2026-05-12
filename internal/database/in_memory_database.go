package database

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/isichei/recipe-book/internal/recipes"
)

type InMemDatabase struct {
	data       map[string]recipes.RecipeMetadata
	sc         SearchCache
	fileGetter MarkdownFileGetter
}

func NewTestDatabaseFromDir(dir string) InMemDatabase {

	fileGetter := LocalMarkdownFileGetter{dir}

	data := make(map[string]recipes.RecipeMetadata)
	sc := make(SearchCache)
	for _, file := range fileGetter.files() {
		uid, fullRecipe := fileGetter.getRecipeFromFilePath(file)
		d, _, _ := strings.Cut(filepath.Base(fullRecipe.Description), ".")

		data[uid] = recipes.RecipeMetadata{Uid: uid, Title: fullRecipe.Title, Description: d}
		for _, word := range strings.Split(strings.ToLower(fullRecipe.Title), " ") {
			sc.Add(word, uid)
		}
		for _, ingredient := range fullRecipe.Ingredients {
			sc.Add(strings.ToLower(ingredient.Name), uid)
		}
	}
	fmt.Printf("Total number of recipes read: %d\n", len(data))
	fmt.Printf("Total size of recipe cache: %d\n", len(sc))
	fmt.Printf("Using localFileGetter with dir: %s\n", dir)
	return InMemDatabase{data, sc, fileGetter}
}

// Rubbish search to fill in for a proper search query later
func (db InMemDatabase) SearchRecipes(text string) []recipes.RecipeMetadata {
	var filtered []recipes.RecipeMetadata

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

func (db InMemDatabase) GetRecipeMetadata(recipeUid string) recipes.RecipeMetadata {
	for _, d := range db.data {
		if d.Uid == recipeUid {
			return d
		}
	}
	return recipes.RecipeMetadata{}
}

func (db InMemDatabase) GetRecipe(recipeUid string) recipes.Recipe {
	return db.fileGetter.getRecipe(recipeUid)
}

func (db InMemDatabase) AddRecipe(recipeUid string, r recipes.Recipe) error {
	return errors.New("Not implemented")
}

func (db InMemDatabase) DeleteRecipe(recipeUid string) error {
	// Check if recipe exists
	_, exists := db.data[recipeUid]
	if !exists {
		return fmt.Errorf("recipe %s not found", recipeUid)
	}

	// Delete from data map
	delete(db.data, recipeUid)

	// Clean up search cache - remove the recipeUid from all search terms
	for term, uidSet := range db.sc {
		delete(uidSet, recipeUid)
		// If the set is now empty, remove the term entirely
		if len(uidSet) == 0 {
			delete(db.sc, term)
		}
	}

	// Also need to remove search terms associated with this recipe
	// Get the full recipe to find all its search terms
	fullRecipe := db.fileGetter.getRecipe(recipeUid)
	
	// Remove title search terms
	for _, word := range strings.Split(strings.ToLower(fullRecipe.Title), " ") {
		if uidSet, exists := db.sc[word]; exists {
			delete(uidSet, recipeUid)
			if len(uidSet) == 0 {
				delete(db.sc, word)
			}
		}
	}

	// Remove ingredient search terms
	for _, ingredient := range fullRecipe.Ingredients {
		term := strings.ToLower(ingredient.Name)
		if uidSet, exists := db.sc[term]; exists {
			delete(uidSet, recipeUid)
			if len(uidSet) == 0 {
				delete(db.sc, term)
			}
		}
	}

	return nil
}
