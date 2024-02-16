package main

import (
	"net/http"

	"github.com/isichei/recipe-book/internal/database"
	"github.com/isichei/recipe-book/internal/recipes"
	"github.com/isichei/recipe-book/internal/views"
)
		
// handler for home page
func handlerRoot(db database.RecipeDatabase) http.Handler {
	search_view := views.SearchResults(db.SearchRecipes(""))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		views.Home(search_view).Render(r.Context(), w)
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
		views.OldHome(search_view, text).Render(r.Context(), w)
	})
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

		recipe := recipes.ParseMarkdownFile("./static/recipe_mds/" + recipeUid + ".md")
		views.Recipe(recipe, recipeUid).Render(r.Context(), w)
	})
}
