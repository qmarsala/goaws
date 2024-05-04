package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GoAwsDbRecord struct {
	*gorm.Model
	Key   string `gorm:"key"`
	Value string `gorm:"value"`
}

func ConnectDB() *gorm.DB {
	var host = os.Getenv("POSTGRES_HOST")
	var port = os.Getenv("POSTGRES_PORT")
	var user = os.Getenv("POSTGRES_USER")
	var pass = os.Getenv("POSTGRES_PASS")
	connectionStringTpl := "host=%v user=%v password=%v dbname=postgres port=%v sslmode=disable TimeZone=UTC"
	connectionString := fmt.Sprintf(connectionStringTpl, host, user, pass, port)
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		fmt.Print(":(")
		panic("failed to connect database")
	}
	return db
}
