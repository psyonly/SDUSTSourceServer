/*
本包中集合了若干用于向数据库增删改查的功能函数
目的在于解耦，不使系统直接和数据库打交道
使设计中的一些管理器与DB进行交互，作为中间代理来执行数据库的访问

*/

package gotools

import (
	"Gin/gorm2/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

const CookieAge int = 720 //how long cookie alive

// 用Cookie作为中间件，由于新的Session的设计已经不用此功能
func AuthMiddleWareCookie(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cookie, err := c.Cookie("goWebServer"); err == nil {
			loginU := models.SDUSTUser{}
			db.Find(&loginU, "stu_no = ?", cookie)
			if cookie == loginU.StuNo {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"console Message : ": "登录身份已过期，请重新登录。",
		})
		c.Abort()
		return
	}
}

// 用manager作为代理行使中间件的功能
// 用manager中的provider来与上下文的cookie进行判断，检查当前用户是否合法，是否时间有效
// 从而允许执行后面的过程，否则返回信息要求重登
func AuthMiddleWare(manager *Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exist := manager.CheckSession(c.Request); exist {
			c.Next()
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"console Message : ": "登录身份已过期，请重新登录。",
		})
		c.Abort()
		return
	}
}

// 返回一个根据Cookie中保存的信息从数据库中查询得出的用户信息和对应的权限
// 在早先的版本中数据库直接与Cookie打交道，耦合性强不利于程序的健壮性，现已不用
func CheckUserFromCookie(c *gin.Context, db *gorm.DB) (u models.SDUSTUser, message string) {
	if cookie, err := c.Cookie("goWebServer"); err == nil {
		db.Find(&u, "stu_no = ?", cookie)
		if u.StuNo == cookie {
			message = "permit"
			return
		}
	}
	message = "illegal Account."
	return
}

//
func CheckUserFromSession(ss Session, db *gorm.DB) (u models.SDUSTUser) {
	fmt.Println(PrintLogHead(), ss.Get("accountNum"))
	db.Find(&u, "stu_no = ?", ss.Get("accountNum"))
	return
}

// 凭借UAS的id获取数据库中对应的记录并以实例返回
func CheckUASFromDB(uid string, db *gorm.DB) (uas models.UserAddrSpace) {
	db.Find(&uas, "user_id = ?", uid)
	return
}

func SetUsersCookie(c *gin.Context, u models.SDUSTUser, ServerIP string) {
	c.SetCookie(
		"goWebServer",
		u.StuNo,
		CookieAge,
		"/",
		ServerIP,
		false,
		true)
}
