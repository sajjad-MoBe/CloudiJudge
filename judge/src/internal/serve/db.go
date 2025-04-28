package serve

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var publishedProblemsCount int64
var notPublishedProblemsCount int64

func connectDatabase() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Tehran",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"),
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

	minutesAgo := time.Now().Add(-60 * time.Minute)
	db.Model(&Submission{}).
		Where("updated_at < ? AND status = ?", minutesAgo, "waiting").
		Update("status", "Compilation failed")

	var submissions []Submission
	err = db.Model(&Submission{}).
		Where("updated_at > ? AND status = ?", minutesAgo, "waiting").
		Preload("Problem").
		Find(&submissions).Error

	if err == nil {
		for _, submission := range submissions {
			sendCodeToRun(submission, submission.Problem)
		}
	}
	err = db.Model(&Problem{}).
		Where("is_published = ?", true).
		Count(&publishedProblemsCount).Error
	if err != nil {
		log.Fatalf("Failed to load published problems: %v", err)
		os.Exit(1)
	}
	err = db.Model(&Problem{}).
		Where("is_published = ?", false).
		Count(&notPublishedProblemsCount).Error
	if err != nil {
		log.Fatalf("Failed to load not published problems: %v", err)
		os.Exit(1)
	}
	// log.Println("Database connected and migrated")
}

func CreateAdmin(email string) {
	connectDatabase()
	email = strings.ToLower(email)
	var user User
	if result := db.Where("email = ?", email).First(&user); result.Error == nil {
		user.IsAdmin = true
		db.Save(&user)
		fmt.Println("Admin was created, use its old password.")
		return
	} else {
		password := GenerateRandomToken(8)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("Admin user was not created.")
			return
		}
		user = User{
			Email:    email,
			Password: string(hashedPassword),
			IsAdmin:  true,
		}
		db.Create(&user)
		db.Commit()
		fmt.Println("admin was created, password:", password)
	}
}
