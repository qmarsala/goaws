package main

import (
	"context"
	"fmt"
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
			Where("exception <> 1").
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

		newIndex := calculateHandicapIndex(rounds)
		fmt.Println("new index: ", newIndex)
		if newIndex.Current == currentIndex.Current {
			return map[string]interface{}{}, nil
		}
		if err := db.Model(goaws.HandicapIndex{}).Create(&newIndex).Error; err != nil {
			return nil, err
		}
		return map[string]interface{}{}, nil
	})
}

func calculateHandicapIndex(rounds []goaws.Round) goaws.HandicapIndex {
	return goaws.HandicapIndex{
		Current: (nOutOfTwentyAverage(rounds)) - float32(getIndexAdjustment(len(rounds))),
	}
}

func nOutOfTwentyAverage(rounds []goaws.Round) float32 {
	diffCount := getDiffCountPerRounds(len(rounds))
	scoreDiffs := getScoreDifferentials(rounds)
	sum := float32(0)
	for _, s := range scoreDiffs[:diffCount] {
		sum += s
	}
	return sum / float32(diffCount)
}

func getScoreDifferentials(rounds []goaws.Round) []float32 {
	scoreDifferentials := []float32{}
	for _, r := range rounds {
		scoreDifferentials = append(scoreDifferentials, calculateScoreDifferential(r))
	}
	slices.Sort(scoreDifferentials)
	return scoreDifferentials
}

func calculateScoreDifferential(round goaws.Round) float32 {
	const pcc int = 0 // todo: pcc, not totally sure what it is, though it is 0 most of the time
	fmt.Println("calculating differential: ", round)
	diff := (113 / round.SlopeRating) * (float32(round.AdjustedGrossScore) - round.CourseRating - float32(pcc))
	fmt.Println(diff)
	return diff
}

func getDiffCountPerRounds(roundCount int) int {
	switch {
	case roundCount < 6:
		return 1
	case roundCount < 9:
		return 2
	case roundCount < 12:
		return 3
	case roundCount < 15:
		return 4
	case roundCount < 17:
		return 6
	case roundCount < 19:
		return 7
	default:
		return 8
	}
}

func getIndexAdjustment(roundCount int) int {
	switch {
	case roundCount == 3:
		return 2
	case roundCount == 4 || roundCount == 6:
		return 1
	default:
		return 0
	}
}
