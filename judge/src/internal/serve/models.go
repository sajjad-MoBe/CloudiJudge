package serve

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string    `json:"email" gorm:"unique"`
	Password string    `json:"password"`
	Problems []Problem `gorm:"foreignKey:UserID" json:"problems"`
}

type Problem struct {
	gorm.Model
	Title       string     `gorm:"type:varchar(255);not null" json:"title"`
	Statement   string     `gorm:"type:text;not null" json:"statement"`
	IsPublished bool       `gorm:"default:false" json:"is_published"`
	PublishedAt *time.Time `json:"published_at"`
	TimeLimit   int        `gorm:"default:0" json:"time_limit"`   // in milliseconds
	MemoryLimit float32    `gorm:"default:0" json:"memory_limit"` // in mb
	UserID      uint       `json:"user_id"`
}
