package api

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {

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
func errorResponse(status int, message string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       message,
		StatusCode: status,
		Headers:    map[string]string{"Content-Type": "text/plain"},
	}
}

// This creates handles the actual lambda request and will deal with errors in this handler, any remaining uncaught errors
// should be caught by the lambda itself
func RequestHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	requester, err := RequesterFactory(request)
	if err != nil {
		return errorResponse(405, err.Error()), nil
	}

	resp, err := requester.Do(ctx, request)
	if err != nil {
		return errorResponse(500, err.Error()), nil
	} else {
		return resp, nil
	}
}
