package models

import (
	"github.com/jinzhu/gorm"
	"strconv"
)

type UserAddrSpace struct {
	gorm.Model
	UserID       string `gorm:"type:varchar(12);unique_index"`
	UserAddr     string `gorm:"type:varchar(50)"`
	CurrentSpace int64  `gorm:"type:BIGINT"`
	MAXSpace     int64  `gorm:"type:BIGINT"`
}

func (ua *UserAddrSpace) Set(uid string, uAddr string, cs, ms int64) {
	ua.UserID = uid
	ua.UserAddr = uAddr
	ua.CurrentSpace = cs
	ua.MAXSpace = ms
}

func (ua *UserAddrSpace) GetUseRate() (str string) {
	cur := int(ua.CurrentSpace) / 1024
	max := int(ua.MAXSpace) / 1024

	return strconv.Itoa(cur) + "MB / " + strconv.Itoa(max) + "MB"
}
