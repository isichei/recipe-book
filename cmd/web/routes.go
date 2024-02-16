package main

import "net/http"

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()
	mux.Handle("/", redirectOldBrowser(handlerRoot(app.db)))
	mux.Handle("/old", handlerOldRoot(app.db))
	mux.Handle("/search-recipes", handleSearchRecipes(app.db))
	mux.Handle("/view-recipe", viewRecipe(app.db))

	// serve static folder
	static_fs := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", static_fs))

	// Do some typing then add some middleware
	var handler http.Handler = mux
	handler = getOnly(handler)
	return logRequest(handler, app.logger)

}
