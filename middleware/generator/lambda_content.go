package generator

const LAMBDA_CONTENT = `// Code auto generated; PLEASE EDIT.
package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// PLEASE EDIT THIS Handler function
func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 501,
		Body:       "Not Implemented",
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
`