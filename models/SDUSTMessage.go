package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type UsersMessage struct {
	gorm.Model
	Sender   string `gorm:"type:varchar(12)"`
	Receiver string `gorm:"type:varchar(12)"`
	Content  string `gorm:"type:varchar(500)"`
	SendTime time.Time
	READ     bool `gorm:"type:tinyint(1)"`
}
