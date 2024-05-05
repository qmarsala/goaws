package goaws

import "gorm.io/gorm"

type HandicapIndex struct {
	*gorm.Model
	Value float32
}

type Round struct {
	*gorm.Model
	CourseName         string
	CourseRating       float32
	SlopeRating        float32
	HolesPlayed        int
	Score              int
	AdjustedGrossScore int
}
