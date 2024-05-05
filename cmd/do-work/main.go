package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
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
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("Couldn't load default configuration. Have you set up your AWS account?")
	}
	queueUrl := os.Getenv("OUTPUT_QUEUE_URL")
	sqsClient := sqs.NewFromConfig(sdkConfig)

	lambda.Start(func(ctx context.Context, event *GoAwsRequest) (*GoAwsResponse, error) {
		if event == nil {
			return nil, fmt.Errorf("received nil event")
		}
		fmt.Println("work, work.")
		resp := GoAwsResponse{
			Key:   event.Key,
			Value: event.Value + " - I worked on this :)",
		}
		jsonBody, _ := json.Marshal(resp)
		body := string(jsonBody)
		input := &sqs.SendMessageInput{
			QueueUrl:    &queueUrl,
			MessageBody: &body,
		}
		sqsClient.SendMessage(ctx, input)
		return &resp, nil
	})
}
