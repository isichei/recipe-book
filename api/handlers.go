package api

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

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

	log.Printf("Path: %s", request.Path)
	log.Printf("Method: %s", request.HTTPMethod)

	if request.HTTPMethod == "GET" {
		switch request_path := request.Path; {
		case request_path == "/":
			requester = NewHtmlRequester(true, "")
		case strings.HasPrefix(request_path, "/search-recipes"):
			requester = NewHtmlRequester(false, request.QueryStringParameters["text"])
		case strings.HasPrefix(request_path, "/static"):
			requester = TextRequester{}
		case strings.HasPrefix(request_path, "/thumbnails"):
			image_name, _ := strings.CutPrefix(request_path, "/")
			requester = NewImageRequester(image_name)
		default:
			err_body := fmt.Sprintf("Unknown Path: %s", request.Path)
			log.Print(err_body)
			return events.APIGatewayProxyResponse{
				Body:       err_body,
				StatusCode: 404,
				Headers:    request.Headers,
			}, nil
		}
	} else {
		err_body := "Only expect GET requests"
		log.Print(err_body)
		return events.APIGatewayProxyResponse{
			Body:       err_body,
			StatusCode: 405,
			Headers:    request.Headers,
		}, nil
	}

	body, e := requester.RetrieveData()
	if e != nil {
		body_err := e.Error()
		log.Print(body_err)
		return events.APIGatewayProxyResponse{
			Body:       body_err,
			StatusCode: 405,
			Headers:    request.Headers,
		}, nil
	}
	headers := map[string]string{"Content-Type": requester.ContentType()}
	if requester.ContentType() == "image/jpeg" {
		headers["Content-Length"] = fmt.Sprintf("%d", len(body))
		return events.APIGatewayProxyResponse{
			Body:            base64.StdEncoding.EncodeToString(body),
			StatusCode:      200,
			Headers:         headers,
			IsBase64Encoded: true,
		}, nil

	} else {
		return events.APIGatewayProxyResponse{
			Body:       string(body),
			StatusCode: 200,
			Headers:    headers,
		}, nil
	}
}
