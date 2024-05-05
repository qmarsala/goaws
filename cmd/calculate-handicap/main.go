package main

import (
	"context"
	"fmt"
	"math"
	"slices"

	goaws "goaws/internal"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	db := goaws.ConnectDB()
	lambda.Start(func(ctx context.Context, sqsEvent events.SQSEvent) (map[string]interface{}, error) {
		fmt.Println("score posted, recalculating handicap")
		rounds := []goaws.Round{}
		db.Model(goaws.Round{}).Limit(20).Order("CreatedAt DESC").Find(&rounds)

		currentIndex := goaws.HandicapIndex{}
		db.Model(goaws.HandicapIndex{}).Order("CreatedAt DESC").First(currentIndex)

		newIndex := calculateHandicapIndex(rounds)
		if newIndex.Value == currentIndex.Value {
			return map[string]interface{}{}, nil
		}
		if err := db.Model(goaws.HandicapIndex{}).Create(&newIndex).Error; err != nil {
			return nil, err
		}
		return map[string]interface{}{}, nil
	})
}

func calculateScoreDifferential(round goaws.Round) float32 {
	const pcc int = 0 // todo: pcc, not totally sure what it is, though it is 0 most of the time
	return (113 / round.SlopeRating) * (float32(round.AdjustedGrossScore) - round.CourseRating - float32(pcc))
}

func getTopEightScoreDifferentials(rounds []goaws.Round) []float32 {
	topEight := []float32{}
	for _, r := range rounds {
		slices.Sort(topEight)
		slices.Reverse(topEight)
		scoreDifferential := calculateScoreDifferential(r)
		if len(topEight) < 8 {
			topEight = append(topEight, scoreDifferential)
		} else if scoreDifferential < topEight[0] {
			topEight[0] = scoreDifferential
		}
	}
	return topEight
}

func calculateHandicapIndex(rounds []goaws.Round) goaws.HandicapIndex {
	sum := float32(0)
	topEight := getTopEightScoreDifferentials(rounds)
	for _, s := range topEight {
		sum += s
	}
	//todo: apply any safe guards for needed for posted round
	return goaws.HandicapIndex{
		Value: sum / float32(math.Max(1, float64(len(topEight)))),
	}
}
