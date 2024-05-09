package goaws

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type HandicapIndex struct {
	*gorm.Model
	Current float32
	Low     float32
}

type Round struct {
	*gorm.Model
	CourseName            string
	CourseRating          float32
	SlopeRating           float32
	HolesPlayed           int
	Score                 int
	PostedScore           int
	ScoreDifferential     float32
	ExceptionalAdjustment int
	Exceptional           bool
	ThrowAway             bool
}

type DatabaseConnection struct {
	*gorm.DB
}

func ProvideDatabase() (DatabaseConnection, error) {
	if db, err := connectDB(); err == nil {
		return DatabaseConnection{DB: db}, nil
	} else {
		return DatabaseConnection{}, err
	}
}

func connectDB() (*gorm.DB, error) {
	var host = os.Getenv("POSTGRES_HOST")
	var port = os.Getenv("POSTGRES_PORT")
	var user = os.Getenv("POSTGRES_USER")
	var pass = os.Getenv("POSTGRES_PASS")
	connectionStringTpl := "host=%v user=%v password=%v dbname=postgres port=%v sslmode=disable TimeZone=UTC"
	connectionString := fmt.Sprintf(connectionStringTpl, host, user, pass, port)
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		fmt.Print("failed to connect database")
		return nil, err
	}
	return db, nil
}
