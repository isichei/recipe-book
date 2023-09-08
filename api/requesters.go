package api

import (
	"bytes"
	"context"
	"embed"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/isichei/recipe-book/types"
)

// content holds our static web server content. Note path is relative to ./main.go!
// So I set this embedding in main
var StaticResources embed.FS

type Requester interface {
	Do(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

func RequesterFactory(request events.APIGatewayProxyRequest) (Requester, error) {

	log.Printf("Path: %s", request.Path)
	log.Printf("Method: %s", request.HTTPMethod)

	if request.HTTPMethod != "GET" {
		return nil, errors.New(fmt.Sprintf("API only accepts GET requests but got a %s request", request.HTTPMethod))
	}
	var requester Requester
	switch request_path := request.Path; {
	case request_path == "/":
		requester = homeRequester{}
	case strings.HasPrefix(request_path, "/search-recipes"):
		requester = searchRecipesRequester{}
	case strings.HasPrefix(request_path, "/static"):
		requester = staticRequester{}
	case strings.HasPrefix(request_path, "/thumbnails"):
		requester = imageRequester{}
	default:
		return nullRequester{}, errors.New(fmt.Sprintf("Recieved unexpected Path %s", request_path))
	}
	return requester, nil
}

type nullRequester struct{}

func (nullRequester) Do(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{}, nil
}

type homeRequester struct{}

func (homeRequester) Do(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	h := template.Must(template.ParseFS(StaticResources, "templates/home.html", "templates/search_results.html"))
	return htmlTemplateToResponse(h, searchRecipes(""))
}

type searchRecipesRequester struct{}

func (searchRecipesRequester) Do(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	tmpl := template.Must(template.ParseFS(StaticResources, "templates/search_results.html"))
	return htmlTemplateToResponse(tmpl, searchRecipes(request.QueryStringParameters["text"]))
}

type staticRequester struct{}

func (staticRequester) Do(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	b, err := fs.ReadFile(StaticResources, "static/styles.css")
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	headers := map[string]string{"Content-Type": "text/css"}
	return events.APIGatewayProxyResponse{
		Body:       string(b),
		StatusCode: 200,
		Headers:    headers,
	}, nil
}

type imageRequester struct{}

func (imageRequester) Do(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	image_name, _ := strings.CutPrefix(request.Path, "/")
	b, err := fs.ReadFile(StaticResources, image_name)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return imageToResponse(b), nil
}

func htmlTemplateToResponse(tmpl *template.Template, data any) (events.APIGatewayProxyResponse, error) {
	var html bytes.Buffer

	e := tmpl.Execute(&html, data)

	if e != nil {
		return events.APIGatewayProxyResponse{}, e
	}
	headers := map[string]string{"Content-Type": "text/html"}
	return events.APIGatewayProxyResponse{
		Body:       html.String(),
		StatusCode: 200,
		Headers:    headers,
	}, nil
}

func imageToResponse(body []byte) events.APIGatewayProxyResponse {
	headers := map[string]string{"Content-Type": "image/jpeg", "Content-Length": fmt.Sprintf("%d", len(body))}
	return events.APIGatewayProxyResponse{
		Body:            base64.StdEncoding.EncodeToString(body),
		StatusCode:      200,
		Headers:         headers,
		IsBase64Encoded: true,
	}
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
