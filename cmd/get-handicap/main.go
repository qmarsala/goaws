// this will calculate your effective handicap at a particular course
package main

import (
	"context"
	"fmt"

	"goaws/di"
	goaws "goaws/internal"

	"github.com/aws/aws-lambda-go/lambda"
)

type GetCourseHandicapRequest struct {
	SlopeRating float32
}

func main() {
	db, err := di.InitializeDatabase()
	if err != nil {
		panic(err)
	}

	lambda.Start(func(ctx context.Context, event *GetCourseHandicapRequest) (map[string]interface{}, error) {
		if !(event.SlopeRating >= 55 || event.SlopeRating <= 155) {
			return nil, fmt.Errorf("invalid slope rating")
		}

		currentIndex := goaws.HandicapIndex{}
		if err := db.Model(goaws.HandicapIndex{}).
			Order("created_at desc").
			First(&currentIndex).Error; err != nil {
			fmt.Println("Error getting handicap index: ", err)
			return nil, fmt.Errorf("unable to retrieve index")
		}

		courseHandicap := (currentIndex.Current * event.SlopeRating) / 113
		return map[string]interface{}{
			"CourseHandicap": courseHandicap,
		}, nil
	})
}
