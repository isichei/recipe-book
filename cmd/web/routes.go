package main

import (
	"net/http"
)

const ENABLE_AUTH = false

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()
	mux.Handle("GET /", redirectOldBrowser(handlerRoot(app.db)))
	mux.Handle("GET /favicon.ico", http.NotFoundHandler())
	mux.Handle("GET /old", handlerOldRoot(app.db))
	mux.Handle("GET /search-recipes", handleSearchRecipes(app.db))
	mux.Handle("GET /view-recipe", viewRecipe(app.db))

	if app.enableWrite {
		mux.Handle("GET /add-recipe", addRecipe(app.db))
	}

	// serve static folder, either as embedded FS or local FS
	mux.Handle("GET /static/", staticFileServer(app.staticFolder))

	// All requests
	if ENABLE_AUTH {
		return logRequest(basicAuth(mux, app.creds), app.logger)
	} else {
		return logRequest(mux, app.logger)
	}
}
