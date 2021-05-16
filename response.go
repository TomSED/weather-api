package weatherapi

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func internalServerError() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       "Something has gone wrong",
	}
}

func badRequest(message string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       message,
	}
}

func success(body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       body,
	}
}
