package serve

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func connectDatabase() {
	dsn := fmt.Sprintf(
		"host=localhost user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Tehran",
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"),
	)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		os.Exit(1)
	}

	// Automatically migrate the User model
	err = db.AutoMigrate(&User{}, &Problem{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
		os.Exit(1)
	}

	log.Println("Database connected and migrated")
}
