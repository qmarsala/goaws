package main

import (
	"context"
	"encoding/json"
	"fmt"

	database "goaws/internal"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type GoAwsRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func main() {
	db := database.ConnectDB()
	lambda.Start(func(ctx context.Context, sqsEvent events.SQSEvent) (map[string]interface{}, error) {
		batchItemFailures := []map[string]interface{}{}
		fmt.Println("reading output")
		for _, message := range sqsEvent.Records {
			fmt.Println(message.Body)

			var val GoAwsRequest
			json.Unmarshal([]byte(message.Body), &val)
			trx := db.Model(database.GoAwsDbRecord{}).Create(&database.GoAwsDbRecord{Key: val.Key, Value: val.Value})
			if trx.RowsAffected < 1 {
				batchItemFailures = append(batchItemFailures, map[string]interface{}{"itemIdentifier": message.MessageId})
			}
		}

		sqsBatchResponse := map[string]interface{}{
			"batchItemFailures": batchItemFailures,
		}
		return sqsBatchResponse, nil
	})
}
