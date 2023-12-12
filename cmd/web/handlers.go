package main

import (
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
