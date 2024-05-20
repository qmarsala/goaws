package goaws

import (
	"fmt"
	"os"
	"slices"

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

type DatabaseConnectionConfig struct {
	host, port, user, pass string
}

func ProvideConfig() (DatabaseConnectionConfig, error) {
	var host = os.Getenv("POSTGRES_HOST")
	var port = os.Getenv("POSTGRES_PORT")
	var user = os.Getenv("POSTGRES_USER")
	var pass = os.Getenv("POSTGRES_PASS")
	if slices.Contains([]string{host, port, user, pass}, "") {
		return DatabaseConnectionConfig{}, fmt.Errorf("unable to reade config values")
	}
	return DatabaseConnectionConfig{
		host: host,
		port: port,
		user: user,
		pass: pass,
	}, nil
}

func ProvideDatabase(config DatabaseConnectionConfig) (DatabaseConnection, error) {
	if db, err := connectDB(config); err == nil {
		return DatabaseConnection{DB: db}, nil
	} else {
		return DatabaseConnection{}, err
	}
}

func connectDB(config DatabaseConnectionConfig) (*gorm.DB, error) {
	connectionStringTpl := "host=%v user=%v password=%v dbname=postgres port=%v sslmode=disable TimeZone=UTC"
	connectionString := fmt.Sprintf(connectionStringTpl, config.host, config.user, config.pass, config.port)
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		fmt.Print("failed to connect database")
		return nil, err
	}
	return db, nil
}
