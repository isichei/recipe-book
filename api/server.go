package api

import (
	"github.com/isichei/recipe-book/database"
	"html/template"
	"log"
	"net/http"
)

// Filters the slice of recipe metadata based on text and returns
// the filtered slice
type Server struct {
	listenerAddr string
	db           database.RecipeDatabase
}

func NewServer(listenerAddr string) *Server {

	s := Server{listenerAddr, &database.TestDatabase{}}
	return &s
}

func (s *Server) Start() error {
	log.Println("Starting Recipe App...")

	// routes for server
	http.HandleFunc("/", s.getOnly(s.handlerRoot()))
	http.HandleFunc("/search-recipes", s.getOnly(s.handleSearchRecipes()))

	// serve static folder
	static_fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", static_fs))

	// serve thumbnails folder
	thumbnails_fs := http.FileServer(http.Dir("./thumbnails"))
	http.Handle("/thumbnails/", http.StripPrefix("/thumbnails/", thumbnails_fs))

	return http.ListenAndServe(s.listenerAddr, nil)
}

// middle where to check if is a GetRequest
func (s *Server) getOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		} else {
			h(w, r)
		}
	}
}

// handler for home page
func (s *Server) handlerRoot() http.HandlerFunc {

	tmpl := template.Must(template.ParseFiles("templates/home.html", "templates/search_results.html"))
	params := s.db.SearchRecipes("")

	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.String())
		tmpl.Execute(w, params)
	}

}

// handler for the search recipe
func (s *Server) handleSearchRecipes() http.HandlerFunc {

	tmpl := template.Must(template.ParseFiles("templates/search_results.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.String())
		// Parse the form data to retrieve the parameter value
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusInternalServerError)
			return
		}

		recipeMetadata := s.db.SearchRecipes(r.Form.Get("text"))
		tmpl.Execute(w, recipeMetadata)
	}
}
