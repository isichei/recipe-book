package main

import "net/http"

func (app *application) routes() *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.getOnly(app.handlerRoot()))
	mux.HandleFunc("/search-recipes", app.getOnly(app.handleSearchRecipes()))

	// serve static folder
	static_fs := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", static_fs))

	return mux
}
