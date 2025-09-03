package database

import (
	"aggregator/models"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofrs/uuid/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase() error {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", host, port, user, password, dbname)
	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	err = DB.AutoMigrate(&models.UserSub{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	var count int64
	DB.Model(&models.UserSub{}).Count(&count)
	if count == 0 {
		newSub := []models.UserSub{
			{
				ServiceName: "Yandex Plus",
				UserID:      GetUUID(),
				Cost:        400,
				StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:     time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC), // Чтобы не было пусто, подписка на месяц
			},
			{
				ServiceName: "Kinopoisk Premium",
				UserID:      GetUUID(),
				Cost:        500,
				StartDate:   time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC),
				EndDate:     time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC), // А эта на 2 месяца
			},
		}
		if err := DB.Create(&newSub).Error; err != nil {
			log.Printf("Failed to insert initial data: %v", err)
		}
	}

	return nil
}

func GetUUID() uuid.UUID {
	userID, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
	}

	return userID
}
