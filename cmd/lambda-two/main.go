package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type GoAwsOtherRequest struct {
	Name string `json:"name"`
}

type GoAwsOtherResponse struct {
	Message string `json:"message"`
}

func HandleRequest(ctx context.Context, event *GoAwsOtherRequest) (*GoAwsOtherResponse, error) {
	if event == nil {
		return nil, fmt.Errorf("received nil event")
	}
	message := fmt.Sprintf("Hello %s!", event.Name)
	return &GoAwsOtherResponse{Message: message}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
