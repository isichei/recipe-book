package api

import (
	"bytes"
	"context"
	"html/template"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/isichei/recipe-book/types"
)

type SimpleResponse struct {
	Body       string
	StatusCode int
}

// Testing handler for debugging
func TestRequestHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	default_body := "Unknown request: " + request.HTTPMethod + "with path: " + request.Path
	default_status := 404
	sr := SimpleResponse{Body: default_body, StatusCode: default_status}

	if request.HTTPMethod == "GET" {
		switch request.Path {
		case "/":
			sr.Body = "Root hit"
			sr.StatusCode = 200
		case "/search-recipes":
			text := request.QueryStringParameters["text"]
			sr.Body = "Search response hit with query: " + text
			sr.StatusCode = 200
		default:
			sr.Body = "Unknown path: " + request.Path
			sr.StatusCode = 404
		}
	}

	return events.APIGatewayProxyResponse{Body: sr.Body, StatusCode: sr.StatusCode}, nil
}

// Real handler for recipe App
func RecipeRequestHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var htmlBody bytes.Buffer
	status := 500

	if request.HTTPMethod == "GET" {
		switch request.Path {
		case "/":
			tmpl := template.Must(template.ParseFiles("templates/home.html", "templates/search_results.html"))
			tmpl.Execute(&htmlBody, searchRecipes(""))
			status = 200
		case "/search-recipes":

			tmpl := template.Must(template.ParseFiles("templates/search_results.html"))
			tmpl.Execute(&htmlBody, searchRecipes(request.QueryStringParameters["text"]))
			status = 200
		default:
			htmlBody.WriteString("Unknown path: " + request.Path)
			status = 404
		}
	}

	return events.APIGatewayProxyResponse{Body: htmlBody.String(), StatusCode: status}, nil
}

// Todo move stuff around once lambdas are working as this duplicates storage package
func searchRecipes(text string) []types.RecipeMetadata {

	data := []types.RecipeMetadata{
		{
			Uid:         "chicken-dhansak-recipe",
			Title:       "Chicken Dhansak",
			Description: "A chicken dhansak recipe from BBC good foods",
		},
		{
			Uid:         "christmas-roast-potatoes",
			Title:       "Jamie Oliver Roast Potatoes",
			Description: "A jamie oliver roast potato recipe usually used at Christmas",
		},
	}

	if text == "" {
		return data
	} else {
		var filtered []types.RecipeMetadata

		for _, recipe := range data {
			if strings.Contains(strings.ToLower(recipe.Description), strings.ToLower(text)) {
				filtered = append(filtered, recipe)
			}
		}
		return filtered
	}
}
