package main

import (
	"log"
	"net/http"
	"runtime/debug"
	"strings"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {

	log.Printf("ERROR: on %s %s. With trace: %s", r.Method, r.URL.RequestURI(), string(debug.Stack()))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Middleware to check if the browser works with htmx, or more specifically is it my super old iPad
func (app *application) htmxEnabled(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userAgent := r.Header["User-Agent"][0]
		if strings.Contains(userAgent, "iPad") && strings.Contains(userAgent, "9_3_5") {
			
		}
}
}

