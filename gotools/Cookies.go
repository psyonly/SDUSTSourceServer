package gotools

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

const CookieAge int = 720 //how long cookie alive

func AuthMiddleWare(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cookie, err := c.Cookie("goWebServer"); err == nil {
			loginU := SDUSTUser{}
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

//返回一个根据Cookie中保存的信息从数据库中查询得出的用户信息和对应的权限
func CheckUserFromCookie(c *gin.Context, db *gorm.DB) (u SDUSTUser, message string) {
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

func SetUsersCookie(c *gin.Context, u SDUSTUser, ServerIP string) {
	c.SetCookie(
		"goWebServer",
		u.StuNo,
		CookieAge,
		"/",
		ServerIP,
		false,
		true)
}
