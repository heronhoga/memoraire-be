package config

import (
	"fmt"
	"os"

	"github.com/heronhoga/memoraire-be/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DatabaseInit() {
	_ = godotenv.Load()

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok", dbHost, dbUser, dbPassword, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		return
	}

	DB = db
	fmt.Println("Database connection established successfully")

	if DB != nil {
		fmt.Println("migrating models..")
		errMigrate := DB.AutoMigrate(&models.User{}, &models.Memo{})
		if errMigrate != nil {
			fmt.Println("Error migrating models to database")
		} else {
			fmt.Println("Models successfully migrated")
		}
	}

}