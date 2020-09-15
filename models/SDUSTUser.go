package models

import "github.com/jinzhu/gorm"

// 网站系统的用户结构体
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

// 用于检查当前用户的信息完成情况
func (u *SDUSTUser) CheckUsersIFMFinish() (percent float32) {
	all := 0
	if u.Age != 0 {
		all++
	}
	if u.Email != "" {
		all++
	}
	if u.DomNo != "" {
		all++
	}
	if u.StuNo != "" {
		all++
	}
	percent = float32(all*100) / 4
	return
}
