package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()
	mux.Handle("/", getOnly(redirectOldBrowser(handlerRoot(app.db))))
	mux.Handle("/old", getOnly(handlerOldRoot(app.db)))
	mux.Handle("/search-recipes", getOnly(handleSearchRecipes(app.db)))
	mux.Handle("/view-recipe", getOnly(viewRecipe(app.db)))
	
	if app.enableWrite {
		mux.Handle("/add-recipe", addRecipe(app.db))
	}

	// serve static folder, either as embedded FS or local FS
	mux.Handle("/static/", staticFileServer(app.staticFolder))

	// Do some typing then add some middleware
	return logRequest(mux, app.logger)

}
