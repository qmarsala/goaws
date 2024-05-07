package goaws

import (
	"slices"
)

func CalculateHandicapIndex(rounds []Round) float32 {
	return CalculateNOutOfTwentyAverage(rounds) - float32(getIndexAdjustment(len(rounds)))
}

func CalculateNOutOfTwentyAverage(rounds []Round) float32 {
	diffCount := getDiffCountPerRounds(len(rounds))
	scoreDiffs := getScoreDifferentials(rounds)
	sum := float32(0)
	for _, s := range scoreDiffs[:diffCount] {
		sum += s
	}
	return sum / float32(diffCount)
}

func getScoreDifferentials(rounds []Round) []float32 {
	scoreDifferentials := []float32{}
	for _, r := range rounds {
		scoreDifferentials = append(scoreDifferentials, calculateScoreDifferential(r))
	}
	slices.Sort(scoreDifferentials)
	return scoreDifferentials
}

func calculateScoreDifferential(round Round) float32 {
	const pcc int = 0 // todo: pcc, not totally sure what it is, though it is 0 most of the time
	diff := (113 / round.SlopeRating) * (float32(round.AdjustedGrossScore) - round.CourseRating - float32(pcc))
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
