package api

import (
	"bytes"
	"embed"
	"html/template"
	"io/fs"
	"strings"

	"github.com/isichei/recipe-book/types"
)

// content holds our static web server content. Note path is relative to ./main.go!
// So I set this embedding in main
var StaticResources embed.FS

type Requester interface {
	RetrieveData() ([]byte, error)
	ContentType() string
}

type HtmlRequester struct {
	tmpl       *template.Template
	searchText string
}

func NewHtmlRequester(is_root bool, text string) HtmlRequester {
	var h HtmlRequester

	if is_root {
		h = HtmlRequester{
			template.Must(template.ParseFS(StaticResources, "templates/home.html", "templates/search_results.html")),
			"",
		}
	} else {
		h = HtmlRequester{
			template.Must(template.ParseFS(StaticResources, "templates/search_results.html")),
			text,
		}
	}
	return h
}

func (h HtmlRequester) RetrieveData() ([]byte, error) {
	var html bytes.Buffer

	e := h.tmpl.Execute(&html, searchRecipes(h.searchText))
	return html.Bytes(), e
}

func (h HtmlRequester) ContentType() string {
	return "text/html"
}

type ImageRequester struct {
	imagePath string
}

func NewImageRequester(imagePath string) ImageRequester {
	return ImageRequester{imagePath}
}

func (ir ImageRequester) RetrieveData() ([]byte, error) {
	return fs.ReadFile(StaticResources, ir.imagePath)
}

func (ir ImageRequester) ContentType() string {
	return "image/jpeg"
}

type TextRequester struct{}

func (tr TextRequester) RetrieveData() ([]byte, error) {
	return fs.ReadFile(StaticResources, "static/styles.css")
}

func (tr TextRequester) ContentType() string {
	return "text/css"
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
