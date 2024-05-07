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

//todo: support multiple players

type EnterRoundRequest struct {
	CourseName         string  `json:"courseName"`
	CourseRating       float32 `json:"courseRating"`
	SlopeRating        float32 `json:"slopeRating"`
	HolesPlayed        int     `json:"holesPlayed"`
	AdjustedGrossScore int     `json:"adjustedScore"`
	Score              int     `json:"score"`
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

		if event.AdjustedGrossScore < 1 {
			return nil, fmt.Errorf("no score")
		}

		currentIndex := goaws.HandicapIndex{}
		if err := db.Model(goaws.HandicapIndex{}).Order("created_at DESC").First(&currentIndex).Error; err != nil && err.Error() != "record not found" {
			return nil, err
		}

		roundHistory := []goaws.Round{}
		if err := db.Model(goaws.Round{}).Order("created_at DESC").Limit(20).Find(&roundHistory).Error; err != nil {
			return nil, err
		}

		newRound := createNewRound(event, currentIndex, roundHistory)
		if err := db.Model(goaws.Round{}).Create(newRound).Error; err != nil {
			return nil, err
		}

		publish(ctx, *newRound)
		return &EnterRoundResponse{
			Message: fmt.Sprintf("Round %v Posted.", newRound.ID),
		}, nil
	})
}

func createNewRound(event *EnterRoundRequest, currentIndex goaws.HandicapIndex, roundHistory []goaws.Round) *goaws.Round {
	newRound := goaws.Round{
		CourseName:         event.CourseName,
		CourseRating:       event.CourseRating,
		SlopeRating:        event.SlopeRating,
		HolesPlayed:        event.HolesPlayed,
		Score:              event.Score,
		AdjustedGrossScore: event.AdjustedGrossScore,
	}
	if currentIndex.Model != nil {
		//todo: adjusted gross score is not relative to par, so this logic isn't correct yet
		newRound.Exceptional = newRound.AdjustedGrossScore < (int(currentIndex.Current) - 7)
		if len(roundHistory) > 19 {
			newRound.ThrowAway = (goaws.CalculateNOutOfTwentyAverage(roundHistory) + 3) > currentIndex.Low
		}
	}
	return &newRound
}

func createPublisher() func(ctx context.Context, newRound goaws.Round) {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("Couldn't load default configuration. Have you set up your AWS account?")
	}
	queueUrl := os.Getenv("OUTPUT_QUEUE_URL")
	sqsClient := sqs.NewFromConfig(sdkConfig)
	return func(ctx context.Context, newRound goaws.Round) {
		fmt.Println("publishing new round")
		jsonBody, _ := json.Marshal(newRound)
		body := string(jsonBody)
		input := &sqs.SendMessageInput{
			QueueUrl:    &queueUrl,
			MessageBody: &body,
		}
		if _, err := sqsClient.SendMessage(ctx, input); err != nil {
			fmt.Println("Unable to publish: ", err)
		}
	}
}
