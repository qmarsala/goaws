package main

import (
	"context"
	"fmt"
	database "goaws/internal"

	"github.com/aws/aws-lambda-go/lambda"
	"gorm.io/gorm"
)

type GoAwsRequest struct {
	Key string `json:"key"`
}

type GoAwsResponse struct {
	Value string `json:"value"`
}

func main() {
	db := database.ConnectDB()
	lambda.Start(func(ctx context.Context, event *GoAwsRequest) (*GoAwsResponse, error) {
		if event == nil {
			return nil, fmt.Errorf("received nil event")
		}
		res := doWork(db, event.Key)
		//todo: 404 when not found?
		return &GoAwsResponse{Value: res.Value}, nil
	})
}

func doWork(db *gorm.DB, key string) database.GoAwsDbRecord {
	res := database.GoAwsDbRecord{}
	db.Model(res).Where("key = ?", key).First(&res)
	return res
}
