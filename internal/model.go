package goaws

import "gorm.io/gorm"

type HandicapIndex struct {
	*gorm.Model
	Current float32
	Low     float32
}

type Round struct {
	*gorm.Model
	CourseName         string
	CourseRating       float32
	SlopeRating        float32
	HolesPlayed        int
	Score              int
	AdjustedGrossScore int
	Exceptional        bool
}
