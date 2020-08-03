package gotools

import "github.com/jinzhu/gorm"

type SDUSTUser struct {
	gorm.Model
	Name     string
	Password string `gorm:"type:varchar(18)"`
	Age      int
	StuNo    string `gorm:"type:varchar(12);unique_index"`
	Email    string `gorm:"type:varchar(100);unique_index"`
	DomNo    string `gorm:"type:varchar(10)"`
	SID      int    `gorm:"AUTO_INCREMENT"`
}
