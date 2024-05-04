package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type GoAwsRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type GoAwsResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func main() {
	lambda.Start(func(ctx context.Context, event *GoAwsRequest) (*GoAwsResponse, error) {
		if event == nil {
			return nil, fmt.Errorf("received nil event")
		}
		fmt.Println("work, work.")
		return &GoAwsResponse{
			Key:   event.Key,
			Value: event.Value + " - I worked on this :)",
		}, nil
	})
}
