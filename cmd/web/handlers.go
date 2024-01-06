package main

import (
	"github.com/isichei/recipe-book/internal/recipes"
	"github.com/isichei/recipe-book/internal/views"
	"log"
	"net/http"
)

// handler for home page
func (app *application) handlerRoot() http.HandlerFunc {

	search_view := views.SearchResults(app.db.SearchRecipes(""))
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.String())
		views.Home(search_view).Render(r.Context(), w)
	}
}

// handler for the search recipe
func (app *application) handleSearchRecipes() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.String())
		// Parse the form data to retrieve the parameter value
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusInternalServerError)
			return
		}

		recipeMetadata := app.db.SearchRecipes(r.Form.Get("text"))
		views.SearchResults(recipeMetadata).Render(r.Context(), w)
	}
}

// View a recipe
func (app *application) viewRecipe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recipeUid := r.URL.Query().Get("uid")

		recipeMeta := app.db.GetRecipeMetadata(recipeUid)
		log.Printf("%s - matched uid: %s\n", r.URL.String(), recipeMeta.Uid)
		if recipeMeta.Uid == "" {
			http.NotFound(w, r)
			return
		}

		recipe := recipes.ParseMarkdownFile("./static/recipe_mds/" + recipeUid + ".md")
		views.Recipe(recipe, recipeUid).Render(r.Context(), w)
	}
}
