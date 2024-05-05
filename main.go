package main

import (
	database "goaws/internal"
)

func main() {
	db := database.ConnectDB()
	db.AutoMigrate(database.Round{}, database.HandicapIndex{})
}
