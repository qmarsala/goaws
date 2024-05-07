package main

import (
	"context"
	"fmt"

	goaws "goaws/internal"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

//todo: support multiple players

func main() {
	db := goaws.ConnectDB()
	lambda.Start(func(ctx context.Context, sqsEvent events.SQSEvent) (map[string]interface{}, error) {
		fmt.Println("score posted, recalculating handicap")
		rounds := []goaws.Round{}
		if err := db.Model(goaws.Round{}).Limit(20).
			Where("exceptional != true").
			Where("throw_away != true").
			Order("created_at desc").
			Find(&rounds).Error; err != nil {
			fmt.Println("Error getting rounds: ", err)
			return nil, err
		}

		if len(rounds) < 3 {
			fmt.Println("not enough rounds yet")
			return map[string]interface{}{}, nil
		}

		currentIndex := goaws.HandicapIndex{}
		if err := db.Model(goaws.HandicapIndex{}).
			Order("created_at desc").
			First(&currentIndex).Error; err != nil {
			fmt.Println("Error getting handicap index: ", err)
		}

		newIndexValue := goaws.CalculateHandicapIndex(rounds)
		fmt.Println("new index: ", newIndexValue)
		if newIndexValue == currentIndex.Current {
			return map[string]interface{}{}, nil
		}
		newIndex := goaws.HandicapIndex{
			Current: newIndexValue,
			Low:     currentIndex.Low,
		}
		if newIndexValue < currentIndex.Low {
			newIndex.Low = newIndexValue
		}
		if err := db.Model(goaws.HandicapIndex{}).Create(&newIndex).Error; err != nil {
			return nil, err
		}
		return map[string]interface{}{}, nil
	})
}
