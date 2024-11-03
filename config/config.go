package config

import (
	"fmt"
	"log"
	"os"
	"tmr-backend/entity"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	DB *gorm.DB
}

func NewConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	db_username := os.Getenv("DB_USERNAME")
	db_password := os.Getenv("DB_PASSWORD")
	db_host := os.Getenv("DB_HOST")
	db_name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", db_username, db_password, db_host, db_name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error initialize database: %s", err)
	}

	db.AutoMigrate(&entity.Subject{}, &entity.Lab{}, &entity.LabSleepHistory{}, &entity.LabCueHistory{})

	return &Config{
		DB: db,
	}
}
