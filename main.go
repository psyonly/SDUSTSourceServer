package main

import (
	tools "Gin/gorm2/gotools"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"strconv"
)

const CookieAge int = 120 //how long cookie alive

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

func checkUsersIFMFinish(u SDUSTUser) (percent float32) {
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

func test() int {
	return -1
}

func authMiddleWare(db *gorm.DB) gin.HandlerFunc {
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

func checkUserFromCookie(c *gin.Context, db *gorm.DB) (u SDUSTUser, message string) {
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

func main() {
	db, err := gorm.Open("mysql", "root:yourSQLpassword@/db1?charset=utf8mb4&parseTime=true&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := gin.Default()

	router.Static("staticFile", "./statics")
	router.StaticFS("/download", http.Dir("Files"))
	router.StaticFS("/downloadMovies", http.Dir("D:\\example"))

	db.AutoMigrate(&SDUSTUser{})

	//router.LoadHTMLFiles("./index.html", "./userSpace.html", "./login.html")
	router.LoadHTMLGlob("templates/*") //load all files if has child dir use like "templates/**/*"

	router.GET("/indexPage", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/", func(c *gin.Context) {
		_, msg := checkUserFromCookie(c, db)
		if msg == "permit" {
			//c.Request.URL.Path = "/userSpace"
			//router.HandleContext(c)
			c.Redirect(http.StatusTemporaryRedirect, "/userSpace")
			return
		}
		c.HTML(http.StatusOK, "login.html", nil)
	})

	router.GET("/loginPage", func(c *gin.Context) { //the same as path "/"
		c.HTML(http.StatusOK, "login.html", nil)
	})

	router.GET("/uploadPage", authMiddleWare(db), func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", nil)
	})

	router.GET("/userSpace", authMiddleWare(db), func(c *gin.Context) {
		u, msg := checkUserFromCookie(c, db)
		fmt.Println(msg)
		c.HTML(http.StatusOK, "userSpace.html", gin.H{
			"finish": checkUsersIFMFinish(u),
			"name":   u.Name + " 欢迎回来",
		})
	})

	router.GET("/downloadCenter", authMiddleWare(db), func(c *gin.Context) {
		c.HTML(http.StatusOK, "downloadCenter.html", gin.H{
			"usersUploadFiles": tools.H5trans("<a href=\"/download\">大家上传的文件</a>"),
			"hostFiles":        tools.H5trans("<a href=\"/downloadMovies\">服务器上的资源</a>"),
		})
	})

	router.POST("/sign", func(c *gin.Context) {
		u := SDUSTUser{
			Model:    gorm.Model{},
			Name:     c.PostForm("UserName"),
			Password: c.PostForm("UserPwd"),
			Age:      0,
			StuNo:    c.PostForm("UserStuNo"),
			Email:    c.PostForm("UserEmail"),
			DomNo:    c.PostForm("UserDomNo"),
			SID:      0,
		}
		u.Age, _ = strconv.Atoi(c.PostForm("UserAge"))
		if u.Name == "" || u.Password == "" {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"alert": tools.JsAlert("用户名或者密码不准为空"),
			})
			return
		}
		db.Create(&u)
		c.HTML(http.StatusOK, "userSpace.html", gin.H{
			"name":   u.Name + " 新人",
			"finish": checkUsersIFMFinish(u),
		})
	})

	router.POST("/login", func(c *gin.Context) {
		if c.PostForm("stuNo") == "" {
			c.HTML(http.StatusOK, "login.html", gin.H{
				"alert": tools.JsAlert("学号不为空"),
			})
		} else {
			u := SDUSTUser{}
			db.Find(&u, "stu_no = ?", c.PostForm("stuNo"))
			fmt.Println(u.Password == c.PostForm("pwd"))
			if u.Name != "" && u.Password == c.PostForm("pwd") { //登陆成功刷新Cookie
				c.SetCookie(
					"goWebServer",
					u.StuNo,
					CookieAge,
					"/",
					"127.0.0.1",
					false,
					true)
				c.HTML(http.StatusOK, "userSpace.html", gin.H{
					"name":   u.Name + "! you came back dog!",
					"finish": checkUsersIFMFinish(u),
				})
			} else {
				fmt.Println(u.Name + " login wrong")
				c.HTML(http.StatusOK, "login.html", gin.H{
					"alert": tools.JsAlert("用户名或密码错误，请重新登录！"),
				})
			}
		}
	})

	router.POST("/usersUpload", authMiddleWare(db), func(c *gin.Context) {

		file, err := c.FormFile("usersFile")
		if err != nil {
			fmt.Println("a user's upload has err, please check what has happened.")
			return
		}
		dst := fmt.Sprintf("./Files/UsersUpload/%s", file.Filename)
		c.SaveUploadedFile(file, dst)
		c.HTML(http.StatusOK, "upload.html", gin.H{
			"ifm": tools.H5trans("<p><i>Upload success!</i></p>"),
		})
	})

	router.Run(":9090")
}
