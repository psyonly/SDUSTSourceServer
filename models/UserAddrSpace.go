package models

import (
	"github.com/jinzhu/gorm"
)

type UserAddrSpace struct {
	gorm.Model
	UserID       string `gorm:"type:varchar(12);unique_index"`
	UserAddr     string `gorm:"type:varchar(50)"`
	CurrentSpace uint64 `gorm:"type:BIGINT"`
	MAXSpace     uint64 `gorm:"type:BIGINT"`
}

func (ua *UserAddrSpace) Set(uid string, uAddr string, cs, ms uint64) {
	ua.UserID = uid
	ua.UserAddr = uAddr
	ua.CurrentSpace = cs
	ua.MAXSpace = ms
}
