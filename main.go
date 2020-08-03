package main

import (
	tools "Gin/gorm2/gotools"
	message "Gin/gorm2/messageSys"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net"
	"net/http"
	"strconv"
	"time"
)

var (
	ServerIP string
)

func checkUsersIFMFinish(u tools.SDUSTUser) (percent float32) {
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

func getIP() (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil && ipnet.IP.IsGlobalUnicast() {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("valid local IP not found")
}

func main() {
	db, err := gorm.Open("mysql", "root:psykfysm3tik5*@/db1?charset=utf8mb4&parseTime=true&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.AutoMigrate(&tools.SDUSTUser{})
	db.AutoMigrate(&message.UsersMessage{})

	ServerIP, _ = getIP()
	fmt.Println("Get Server ip is ", ServerIP)

	router := gin.Default()

	router.Static("staticFile", "./statics")
	router.StaticFS("/download", http.Dir("Files"))
	router.StaticFS("/downloadMovies", http.Dir("D:\\Movies"))

	//router.LoadHTMLFiles("./index.html", "./userSpace.html", "./login.html")
	router.LoadHTMLGlob("templates/*") //load all files; if has child dir use like "templates/**/*"

	router.GET("/indexPage", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/", func(c *gin.Context) {
		_, msg := tools.CheckUserFromCookie(c, db)
		if msg == "permit" {
			//c.Request.URL.Path = "/userSpace"
			//router.HandleContext(c)
			c.Redirect(http.StatusTemporaryRedirect, "/userSpace")
			return
		}
		c.HTML(http.StatusOK, "login.html", nil)
	})

	router.GET("/loginPage", func(c *gin.Context) {
		//the same as path "/" 重新登陆接口
		//c.Redirect(http.StatusTemporaryRedirect, "/")
		c.HTML(http.StatusOK, "login.html", nil)
	})

	router.GET("/uploadPage", tools.AuthMiddleWare(db), func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", nil)
	})

	router.GET("/test", func(c *gin.Context) {
		c.HTML(http.StatusOK, "userMessageBox.html", gin.H{
			"message": "test1",
			"time":    "2020",
		})
		/*c.JSON(http.StatusOK, gin.H{
			"message": "test1",
			"time":    "2020",
		})*/
	})

	userSpace := router.Group("/userSpace", tools.AuthMiddleWare(db))
	{
		userSpace.Static("staticFile", "./statics") //路由组内重新设置静态路径

		//用户主界面
		userSpace.GET("", func(c *gin.Context) {
			u, msg := tools.CheckUserFromCookie(c, db)
			fmt.Println(msg)
			c.HTML(http.StatusOK, "userSpace.html", gin.H{
				"finish": checkUsersIFMFinish(u),
				"name":   u.Name + " 欢迎回来",
			})
		})

		//用户消息界面
		userSpace.GET("/message", func(c *gin.Context) {
			u, log := tools.CheckUserFromCookie(c, db)
			if log != "permit" {
				c.JSON(http.StatusOK, gin.H{
					"message": "user is not permit to access to this page.",
				})
				return
			}
			msg := []message.UsersMessage{}
			db.Where("Receiver = ?", u.StuNo).Find(&msg)
			str, _ := json.Marshal(msg) //json序列化对象（此处为Byte）
			//fmt.Println(string(str))
			//c.JSON(http.StatusOK, msg)
			c.HTML(http.StatusOK, "userMessageBox.html", gin.H{ //转换为string类型以模板的形式传递给前端
				"msg": string(str),
			})
		})

		userSpace.POST("/message/sendMSG", func(c *gin.Context) {
			u, log := tools.CheckUserFromCookie(c, db)
			r := tools.SDUSTUser{}
			db.Find(&r, "stu_no = ?", c.PostForm("MSG_receiver")) //从提交的表单中获取接收者信息
			msgR := []message.UsersMessage{}
			db.Where("Receiver = ?", u.StuNo).Find(&msgR)
			str, _ := json.Marshal(msgR)
			alert := "" //查找不到反馈
			if r.Name == "" {
				alert = "用户不存在"
			} else {
				if log != "permit" {
					c.JSON(http.StatusOK, gin.H{
						"message": "user is not permit to access to this page.",
					})
					return
				}
				msgS := message.UsersMessage{
					Model:    gorm.Model{},
					Sender:   u.StuNo,
					Receiver: c.PostForm("MSG_receiver"),
					Content:  c.PostForm("MSG_content"),
					SendTime: time.Now(),
					READ:     false,
				}
				db.Create(&msgS)
			}
			c.HTML(http.StatusOK, "userMessageBox.html", gin.H{
				"msg":   string(str),
				"alert": alert,
			})
		})
	}

	router.GET("/downloadCenter", tools.AuthMiddleWare(db), func(c *gin.Context) {
		c.HTML(http.StatusOK, "downloadCenter.html", gin.H{
			"usersUploadFiles": tools.H5trans("<a href=\"/download\">大家上传的文件</a>"),
			"hostFiles":        tools.H5trans("<a href=\"/downloadMovies\">服务器上的资源</a>"),
		})
	})

	router.POST("/sign", func(c *gin.Context) {
		u := tools.SDUSTUser{
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
		tools.SetUsersCookie(c, u, ServerIP)
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
			u := tools.SDUSTUser{}
			db.Find(&u, "stu_no = ?", c.PostForm("stuNo"))
			fmt.Println(u.Password == c.PostForm("pwd"))
			if u.Name != "" && u.Password == c.PostForm("pwd") { //登陆成功刷新Cookie
				tools.SetUsersCookie(c, u, ServerIP)
				c.HTML(http.StatusOK, "userSpace.html", gin.H{
					"name":   u.Name + "! 欢迎回来!",
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

	router.POST("/usersUpload", tools.AuthMiddleWare(db), func(c *gin.Context) {
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
