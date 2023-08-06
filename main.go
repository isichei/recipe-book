package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/isichei/recipe-book/api"
)

func main() {
	lambda.Start(api.RecipeRequestHandler)
}
