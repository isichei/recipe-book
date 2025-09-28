package main

import (
	"fmt"
	"net/http"

	"github.com/isichei/recipe-book/internal/recipes"
	"github.com/isichei/recipe-book/internal/database"
	"github.com/isichei/recipe-book/cmd/web/views"
)

// handler for home page
func handlerRoot(db database.RecipeDatabase) http.Handler {
	search_view := views.SearchResults(db.SearchRecipes(""))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		views.HomeComposition(search_view, false, "").Render(r.Context(), w)
	})
}

func handlerOldRoot(db database.RecipeDatabase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusInternalServerError)
			return
		}
		text := r.Form.Get("text")
		search_view := views.SearchResults(db.SearchRecipes(text))
		views.HomeComposition(search_view, true, text).Render(r.Context(), w)
	})
}

func staticFileServer(staticPath string) http.Handler {
	static_fs := http.FileServer(http.Dir(staticPath))
	return http.StripPrefix("/static/", static_fs)
}

// handler for the search recipe
func handleSearchRecipes(db database.RecipeDatabase) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse the form data to retrieve the parameter value
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusInternalServerError)
			return
		}

		recipeMetadata := db.SearchRecipes(r.Form.Get("text"))
		views.SearchResults(recipeMetadata).Render(r.Context(), w)
	})
}

// View a recipe
func viewRecipe(db database.RecipeDatabase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recipeUid := r.URL.Query().Get("uid")

		recipeMeta := db.GetRecipeMetadata(recipeUid)
		if recipeMeta.Uid == "" {
			http.NotFound(w, r)
			return
		}

		recipe := db.GetRecipe(recipeUid)
		views.Recipe(recipe, recipeUid).Render(r.Context(), w)
	})
}

// Add a recipe form
func addRecipe(db database.RecipeDatabase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var recipeUid string

			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Could not parse recipe data", http.StatusUnprocessableEntity)
				return
			}

			ingredientItems := r.PostForm["ingredient-item[]"]
			ingredientAmounts := r.PostForm["ingredient-item[]"]

			if len(ingredientItems) != len(ingredientAmounts) {
				http.Error(w, "Ingredients did not match up", http.StatusUnprocessableEntity)
				return
			}
			var ingredients []recipes.Ingredient
			for i, item := range ingredientItems {
				ingredients = append(ingredients, recipes.Ingredient{Name: item, Amount: ingredientAmounts[i]})
			}

			recipe := recipes.Recipe{
				Title: r.Form.Get("title"),
				PrepTime: r.Form.Get("prep-time"),
				CookingTime: r.Form.Get("cook-time"),
				Serves: r.Form.Get("serves"),
				Ingredients: ingredients,
				Method: r.Form["method-step[]"],
				OtherNotes: r.Form.Get("other-notes"),
				Source: r.Form.Get("source"),
			}
			recipeUid = r.Form.Get("uid")
			if recipeUid == "" {
				fmt.Println("Recipe has no uid")
			}
			err = db.AddRecipe(recipeUid, recipe)
			if err == nil {
				fmt.Printf("Recipe %s added\n")
			} else {
				fmt.Printf("Recipe %s errored: %s\n", recipeUid, err)
			}
		} else {
			component := r.URL.Query().Get("component")
			switch component {
			case "":
				views.AddRecipe().Render(r.Context(), w)
			case "method":
				views.CompMethod().Render(r.Context(), w)
			case "ingredient":
				views.CompIngredient().Render(r.Context(), w)
			}
		}
	})
}
