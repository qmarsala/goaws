package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	//db := database.ConnectDB()
	lambda.Start(func(ctx context.Context, sqsEvent events.SQSEvent) (map[string]interface{}, error) {
		batchItemFailures := []map[string]interface{}{}

		for _, message := range sqsEvent.Records {
			fmt.Println(message)
			//trx := db.Model(database.GoAwsDbRecord{}).Create(&database.GoAwsDbRecord{Key: "test", Value: "test"})
			// if trx.RowsAffected < 1 {
			// 	batchItemFailures = append(batchItemFailures, map[string]interface{}{"itemIdentifier": message.MessageId})
			// }
		}

		sqsBatchResponse := map[string]interface{}{
			"batchItemFailures": batchItemFailures,
		}
		return sqsBatchResponse, nil
	})
}
