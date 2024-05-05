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

//todo: support multiple players

func main() {
	db := goaws.ConnectDB()
	lambda.Start(func(ctx context.Context, sqsEvent events.SQSEvent) (map[string]interface{}, error) {
		fmt.Println("score posted, recalculating handicap")
		rounds := []goaws.Round{}
		if err := db.Model(goaws.Round{}).Limit(20).
			Order("created_at desc").
			Find(&rounds).Error; err != nil {
			fmt.Println("Error getting rounds: ", err)
			return nil, err
		}

		if len(rounds) < 8 {
			fmt.Println("not enough rounds yet")
			return map[string]interface{}{}, nil
		}

		currentIndex := goaws.HandicapIndex{}
		if err := db.Model(goaws.HandicapIndex{}).
			Order("created_at desc").
			First(&currentIndex).Error; err != nil {
			fmt.Println("Error getting handicap index: ", err)
		}

		newIndex := calculateHandicapIndex(rounds)
		fmt.Println("new index: ", newIndex)
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
	fmt.Println("calculating differential: ", round)
	diff := (113 / round.SlopeRating) * (float32(round.AdjustedGrossScore) - round.CourseRating - float32(pcc))
	fmt.Println(diff)
	return diff
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
