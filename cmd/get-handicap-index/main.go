// this will retrieve your handicap index
package main

import (
	"context"
	"fmt"

	"goaws/di"
	goaws "goaws/internal"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	db, err := di.InitializeDatabase()
	if err != nil {
		panic(err)
	}

	lambda.Start(func(ctx context.Context, event *interface{}) (goaws.HandicapIndex, error) {
		currentIndex := goaws.HandicapIndex{}
		if err := db.Model(goaws.HandicapIndex{}).
			Order("created_at desc").
			First(&currentIndex).Error; err != nil {
			fmt.Println("Error getting handicap index: ", err)
			return goaws.HandicapIndex{}, fmt.Errorf("unable to retrieve index")
		}
		return currentIndex, nil
	})
}
