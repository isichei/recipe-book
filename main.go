package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

type RecipeMetadata struct {
	Uid string
	Title string
	Description string
}

type rootData struct {
	SearchResults []RecipeMetadata
}

// Filters the slice of recipe metadata based on text and returns
// the filtered slice
func filterReceipeMetadata(recipes []RecipeMetadata, text string) []RecipeMetadata {
	var filtered []RecipeMetadata

	for _, recipe := range recipes {
		if strings.Contains(strings.ToLower(recipe.Description), strings.ToLower(text)) {
			filtered = append(filtered, recipe)
		}
	}

	return filtered
}

// Rubbish search to fill in for a proper search query later
func searchRecipes(text string) []RecipeMetadata {
	recipes := []RecipeMetadata{
		{"chicken-dhansak-recipe", "Chicken Dhansak", "A chicken dhansak recipe from BBC good foods"},
		{"christmas-roast-potatoes", "Jamie Oliver Roast Potatoes", "A jamie oliver roast potato recipe usually used at Christmas"},
	}

	if text == "" {
		return recipes
	} else {
		return filterReceipeMetadata(recipes, text)
	}
}


func main() {
	log.Println("Starting Recipe App...")

	// serve static folder
	static_fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", static_fs))
	
	// serve thumbnails folder
	thumbnails_fs := http.FileServer(http.Dir("./thumbnails"))
	http.Handle("/thumbnails/", http.StripPrefix("/thumbnails/", thumbnails_fs))

	// handler for root
	rootHandler := func(w http.ResponseWriter, r *http.Request) {

		log.Println(r.URL.String())
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		tmpl := template.Must(template.ParseFiles("templates/home.html", "templates/search_results.html"))

		params := rootData{SearchResults: searchRecipes("")}

		tmpl.Execute(w, params)
	}

	// handler for the search recipes
	searchRecipesHandler := func(w http.ResponseWriter, r *http.Request) {

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

		recipeMetadata := searchRecipes(r.Form.Get("text"))

		tmpl := template.Must(template.ParseFiles("templates/search_results.html"))

		tmpl.Execute(w, recipeMetadata)
	}

	// define handlers
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/search-recipes", searchRecipesHandler)

	log.Fatal(http.ListenAndServe(":8000", nil))
}