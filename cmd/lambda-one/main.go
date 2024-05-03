package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type GoAwsRequest struct {
	Name string `json:"name"`
}

type GoAwsResponse struct {
	Message string `json:"message"`
}

func HandleRequest(ctx context.Context, event *GoAwsRequest) (*GoAwsResponse, error) {
	if event == nil {
		return nil, fmt.Errorf("received nil event")
	}
	message := fmt.Sprintf("Hello %s!", event.Name)
	return &GoAwsResponse{Message: message}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
