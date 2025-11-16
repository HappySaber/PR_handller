package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

var DB *sql.DB

type DBConfig struct {
	Host     string
	Port     int
	User     string
	DBName   string
	Password string
}

func buildDBConfig() *DBConfig {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("Invalid DB_PORT: %v", err)
	}

	if isRunningInDocker() {
		return &DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     port,
			User:     os.Getenv("DB_USER"),
			DBName:   os.Getenv("DB_NAME"),
			Password: os.Getenv("DB_PASSWORD"),
		}
	}

	return &DBConfig{
		Host:     os.Getenv("DB_HOST_LOCAL"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PASSWORD"),
	}
}

func (config *DBConfig) dsn() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName,
	)
}

func Init() {
	time.Sleep(5 * time.Second)
	var err error
	dbConfig := buildDBConfig()
	DB, err = sql.Open("postgres", dbConfig.dsn())
	if err != nil {
		log.Fatalf("Error checking database connection: %v", err)
	}

	log.Println("Successfully connected to the database!")
}

func isRunningInDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	return false
}
