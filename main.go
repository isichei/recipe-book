package main

import (
	"embed"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/isichei/recipe-book/api"
)

//go:embed static/* templates/*
var staticResources embed.FS

func main() {

	api.StaticResources = staticResources
	lambda.Start(api.RecipeRequestHandler)
}
