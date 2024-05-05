package main

import (
	"fmt"
	database "goaws/internal"
)

func main() {
	db := database.ConnectDB()
	db.AutoMigrate(database.GoAwsDbRecord{})
	db.Model(database.GoAwsDbRecord{}).Create(&database.GoAwsDbRecord{Key: "hello", Value: "world"})
	results := []database.GoAwsDbRecord{}
	db.Model(database.GoAwsDbRecord{}).Find(&results)
	fmt.Print(results)
}
