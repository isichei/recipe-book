package api

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

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
