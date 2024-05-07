package goaws

func CalculateHandicapIndex(rounds []Round) float32 {
	return CalculateDifferentialAverage(rounds) - float32(getIndexAdjustment(len(rounds)))
}

func CalculateDifferentialAverage(rounds []Round) float32 {
	diffCount := getDiffCountPerRounds(len(rounds))
	adjustedRounds := applyExceptionRoundAdjustments(rounds)
	sum := float32(0)
	for _, r := range adjustedRounds[:diffCount] {
		sum += r.ScoreDifferential
	}
	return sum / float32(diffCount)
}

func CalculateScoreDifferential(round Round) float32 {
	const pcc int = 0 // todo: pcc, not totally sure what it is, though it is 0 most of the time
	diff := (113 / round.SlopeRating) * (float32(round.PostedScore) - round.CourseRating - float32(pcc))
	return diff
}

func applyExceptionRoundAdjustments(roundHistory []Round) []Round {
	totalAdjustment := 0
	for _, r := range roundHistory {
		totalAdjustment = r.ExceptionalAdjustment
	}
	adjustedRounds := []Round{}
	for _, v := range roundHistory {
		adjustedRounds = append(adjustedRounds, Round{
			CourseName:            v.CourseName,
			CourseRating:          v.CourseRating,
			SlopeRating:           v.SlopeRating,
			HolesPlayed:           v.HolesPlayed,
			Score:                 v.Score,
			PostedScore:           v.PostedScore,
			ScoreDifferential:     v.ScoreDifferential - float32(totalAdjustment),
			ExceptionalAdjustment: v.ExceptionalAdjustment,
			Exceptional:           v.Exceptional,
			ThrowAway:             v.ThrowAway,
		})
	}
	return adjustedRounds
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
