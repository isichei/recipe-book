package api

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
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

	var requester Requester

	if request.HTTPMethod == "GET" {
		switch request.Path {
		case "/":
			requester = NewHtmlRequester(true, "")
		case "/search-recipes":
			requester = NewHtmlRequester(false, request.QueryStringParameters["text"])
		case "/static":
			requester = TextRequester{}
		case "/thumbnails":
			requester = ImageRequester{}
		default:
			return events.APIGatewayProxyResponse{
				Body:       fmt.Sprintf("Unknown Path: %s", request.Path),
				StatusCode: 404,
				Headers:    request.Headers,
			}, nil
		}
	} else {
		return events.APIGatewayProxyResponse{
			Body:       "Only expect GET requests",
			StatusCode: 405,
			Headers:    request.Headers,
		}, nil
	}

	body, e := requester.RetrieveData()
	if e != nil {
		return events.APIGatewayProxyResponse{
			Body:       e.Error(),
			StatusCode: 405,
			Headers:    request.Headers,
		}, nil
	}
	headers := map[string]string{"Content-Type": requester.ContentType()}
	if requester.ContentType() == "image/jpeg" {
		headers["Content-Length"] = fmt.Sprintf("%d", len(body))
	}

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
		Headers:    headers,
	}, nil
}
