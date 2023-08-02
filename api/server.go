package api

import (
	"html/template"
	"log"
	"net/http"

	"github.com/isichei/recipe-book/storage"
)

type Server struct {
	port      string
	dataStore storage.Storage
}

func NewServer(listenPort string, dataStore storage.Storage) *Server {
	return &Server{
		port:      listenPort,
		dataStore: dataStore,
	}
}

func (s *Server) Start() error {

	// serve static folder
	static_fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", static_fs))

	// serve thumbnails folder
	thumbnails_fs := http.FileServer(http.Dir("./thumbnails"))
	http.Handle("/thumbnails/", http.StripPrefix("/thumbnails/", thumbnails_fs))

	http.HandleFunc("/", s.rootHandler)
	http.HandleFunc("/search-recipes", s.searchRecipesHandler)

	log.Println("Starting Recipe App at ", s.port)
	return http.ListenAndServe(s.port, nil)
}

// handler for root
func (s *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.String())

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/home.html", "templates/search_results.html"))
	tmpl.Execute(w, s.dataStore.SearchRecipes(""))
}

// handler for the search recipes
func (s *Server) searchRecipesHandler(w http.ResponseWriter, r *http.Request) {

	log.Println(r.URL.String())

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data to retrieve the parameter value
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusInternalServerError)
		return
	}

	recipeMetadata := s.dataStore.SearchRecipes(r.Form.Get("text"))

	tmpl := template.Must(template.ParseFiles("templates/search_results.html"))

	tmpl.Execute(w, recipeMetadata)
}
