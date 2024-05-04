package main

import (
	"fmt"
	database "goaws/internal"
)

func main() {
	db := database.ConnectDB()
	db.AutoMigrate(database.GoAwsDbRecord{})
	db.Model(database.GoAwsDbRecord{}).Create(&database.GoAwsDbRecord{Key: "hello", Value: "world"})
	count := int64(0)
	db.Model(database.GoAwsDbRecord{}).Where("key = ?", "hello").Count(&count)
	fmt.Print(count)
}
