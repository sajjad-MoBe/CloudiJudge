package serve

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email            string    `json:"email" gorm:"unique"`
	Password         string    `json:"password"`
	Problems         []Problem `gorm:"foreignKey:OwnerID" json:"problems"`
	IsAdmin          bool      `gorm:"default:false" json:"isAdmin"`
	AdminCreatedByID uint
}

type Problem struct {
	gorm.Model
	Title       string     `gorm:"type:varchar(50)" json:"title"`
	Statement   string     `gorm:"type:text" json:"statement"`
	IsPublished bool       `gorm:"default:false" json:"is_published"`
	PublishedAt *time.Time `gorm:"default:null" json:"published_at"`
	TimeLimit   int        `gorm:"default:0" json:"time_limit"`   // in milliseconds
	MemoryLimit float32    `gorm:"default:0" json:"memory_limit"` // in mb
	OwnerID     uint       `json:"user_id"`
	Owner       User       `gorm:"foreignKey:OwnerID"`
}

type Submission struct {
	gorm.Model
	Status    string  `gorm:"default:waiting" json:"status"`
	Token     string  `gorm:"type:text" json:"token"`
	OwnerID   uint    `json:"user_id"`
	Owner     User    `gorm:"foreignKey:OwnerID"`
	ProblemID uint    `json:"problem_id"`
	Problem   Problem `gorm:"foreignKey:ProblemID"`
}
