package serve

import (
	"fmt"
	"log"
	"os"
	"time"

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
	err = db.AutoMigrate(&User{}, &Problem{}, &Submission{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
		os.Exit(1)
	}
	tenMinutesAgo := time.Now().Add(-10 * time.Minute)
	db.Model(&Submission{}).
		Where("updated_at < ? AND status == ?", tenMinutesAgo, "waiting").
		Update("status", "Compilation failed")

	log.Println("Database connected and migrated")
}
