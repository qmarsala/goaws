package main

import (
	"context"
	"encoding/json"
	"fmt"
	goaws "goaws/internal"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type EnterRoundRequest struct {
	CourseName   string  `json,required:"courseName"`
	CourseRating float32 `json,required:"courseRating"`
	HolesPlayed  int     `json,required:"holesPlayed"`
	Score        int     `json,required:"score"`
}

type EnterRoundResponse struct {
	Message string `json:"message"`
}

func main() {
	db := goaws.ConnectDB()
	publish := createPublisher()
	lambda.Start(func(ctx context.Context, event *EnterRoundRequest) (*EnterRoundResponse, error) {
		if event == nil {
			return nil, fmt.Errorf("received nil event")
		}

		if !(event.HolesPlayed == 9 || event.HolesPlayed == 18) {
			return nil, fmt.Errorf("only 9 or 18 hole rounds are currently supported")
		}

		newRound := goaws.Round{
			CourseName:   event.CourseName,
			CourseRating: event.CourseRating,
			HolesPlayed:  event.HolesPlayed,
			Score:        event.Score,
		}

		if err := db.Model(goaws.Round{}).Create(&newRound).Error; err != nil {
			return nil, err
		}

		publish(ctx, newRound)
		return &EnterRoundResponse{
			Message: fmt.Sprintf("Round %v Posted.", newRound.ID),
		}, nil
	})
}

func createPublisher() func(ctx context.Context, newRound goaws.Round) {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("Couldn't load default configuration. Have you set up your AWS account?")
	}
	queueUrl := os.Getenv("OUTPUT_QUEUE_URL")
	sqsClient := sqs.NewFromConfig(sdkConfig)
	return func(ctx context.Context, newRound goaws.Round) {
		jsonBody, _ := json.Marshal(newRound)
		body := string(jsonBody)
		input := &sqs.SendMessageInput{
			QueueUrl:    &queueUrl,
			MessageBody: &body,
		}
		sqsClient.SendMessage(ctx, input)
	}
}
