package db

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string
	Nodes []Node `gorm:"foreignkey:UserID"`
	Messages []Message `gorm:"foreignkey:UserID"`

}

type Node struct{
	gorm.Model
	UserId uint
	Name string
	MacAddress string
	ApiKey string
	MessageQueue []Message `gorm:"foreignkey:MessageID"`
}

type Message struct{
	gorm.Model
	UserId uint
	User User `gorm:"foreignkey:UserID"`

	Title string
	Message string
	SentAt time.Time
	Status bool
}
