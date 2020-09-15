package main

import (
	tools "Gin/gorm2/gotools"
	model "Gin/gorm2/models"
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

var SessionManager *tools.Manager
var UASManager *tools.UASManager
var db *gorm.DB

func init() {
	// Session operate
	mp := tools.MapProvider{}
	mp.MPInit()
	tools.Register("MemProvider", &mp)
	SessionManager, _ = tools.NewManager("MemProvider", "goServer", 60)

	// UAS operate
	UASManager = tools.NewUasManager("./UAS", "usr", 20*2<<20)

	// IP operate
	ServerIP, _ = getIP()
	fmt.Println(tools.AddLogHead("Get Server ip is "), ServerIP)

	// DB operate
	var err error
	db, err = gorm.Open("mysql", "root:psykfysm3tik5*@/db1?charset=utf8mb4&parseTime=true&loc=Local")
	if err != nil {
		panic(err)
	}
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

func IndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func Root(c *gin.Context) {
	_, exist := SessionManager.CheckSession(c.Request)
	if exist {
		c.Redirect(http.StatusTemporaryRedirect, "/userSpace")
		return
	}
	c.HTML(http.StatusOK, "login.html", nil)
}

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func UploadPage(c *gin.Context) {
	ss, _ := SessionManager.CheckSession(c.Request)
	c.HTML(http.StatusOK, "upload.html", gin.H{
		"userHeadImageDir": UASManager.GetUASFileDir(ss.Get("accountNum").(string), "HeadImage", db),
	})
}

func UserSpaceRoot(c *gin.Context) {
	ss, _ := SessionManager.CheckSession(c.Request)
	u := tools.CheckUserFromSession(ss, db)
	c.HTML(http.StatusOK, "userSpace.html", gin.H{
		"finish":           u.CheckUsersIFMFinish(),
		"name":             u.Name + " 欢迎回来",
		"userHeadImageDir": UASManager.GetUASFileDir(ss.Get("accountNum").(string), "HeadImage", db),
	})
}

func UserSpaceMessage(c *gin.Context) {
	ss, _ := SessionManager.CheckSession(c.Request)
	u := tools.CheckUserFromSession(ss, db)
	fmt.Println(u)
	msg := []model.UsersMessage{}
	db.Where("Receiver = ?", u.StuNo).Find(&msg)
	str, _ := json.Marshal(msg)                         //json序列化对象（此处为Byte）
	c.HTML(http.StatusOK, "userMessageBox.html", gin.H{ //转换为string类型以模板的形式传递给前端
		"msg":  string(str),
		"user": u.Name,
	})
}

func UserSpaceMessageSendMSGPost(c *gin.Context) {
	ss, _ := SessionManager.CheckSession(c.Request)
	u := tools.CheckUserFromSession(ss, db)
	r := model.SDUSTUser{}
	db.Find(&r, "stu_no = ?", c.PostForm("MSG_receiver")) //从提交的表单中获取接收者信息
	msgR := []model.UsersMessage{}
	db.Where("Receiver = ?", u.StuNo).Find(&msgR) // 从数据库中查找当前用户的消息列表
	str, _ := json.Marshal(msgR)
	alert := "" //查找不到反馈
	if r.Name == "" {
		alert = "用户不存在"
	} else {
		msgS := model.UsersMessage{
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
	//c.Redirect(http.StatusFound, "")
}

func UserSpaceHeadImagePost(c *gin.Context) {
	file, err := c.FormFile("headImage")
	if err != nil {
		fmt.Println(tools.PrintLogHead(), "there is an err happened when upload the image.")
		return
	}
	ss, _ := SessionManager.CheckSession(c.Request)
	uasID := ss.Get("accountNum")
	uas := &model.UserAddrSpace{}
	db.Find(&uas, "user_id = ?", uasID)
	dst := fmt.Sprintf("%s/%s/HeadImage/%s", UASManager.GetFileDir(), uas.UserAddr, "headImage.jpg")
	fmt.Println(tools.PrintLogHead(), "Dir is ", dst)
	c.SaveUploadedFile(file, dst)
	c.Redirect(http.StatusFound, "/userSpace")
}

func DownLoadCenter(c *gin.Context) {
	ss, _ := SessionManager.CheckSession(c.Request)
	c.HTML(http.StatusOK, "downloadCenter.html", gin.H{
		"usersUploadFiles": tools.H5trans("<a href=\"/download\">大家上传的文件</a>"),
		"hostFiles":        tools.H5trans("<a href=\"/downloadMovies\">服务器上的资源</a>"),
		"userHeadImageDir": UASManager.GetUASFileDir(ss.Get("accountNum").(string), "HeadImage", db),
	})
}

func SignPagePost(c *gin.Context) {
	u := model.SDUSTUser{
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
	ss := SessionManager.SessionStart(c.Writer, c.Request)
	ss.Set("accountName", u.Name)
	ss.Set("accountNum", u.StuNo)
	uas := UASManager.InitUserAddrSpace(u.StuNo)
	db.Create(&uas)
	c.HTML(http.StatusOK, "userSpace.html", gin.H{
		"name":   u.Name + " 新人",
		"finish": u.CheckUsersIFMFinish(),
	})
}

func LoginPagePost(c *gin.Context) {
	if c.PostForm("stuNo") == "" {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"alert": tools.JsAlert("学号不为空"),
		})
	} else {
		u := model.SDUSTUser{}
		db.Find(&u, "stu_no = ?", c.PostForm("stuNo"))
		fmt.Println(u.Password == c.PostForm("pwd"))
		if u.Name != "" && u.Password == c.PostForm("pwd") { //登陆成功刷新Cookie/Session
			ss := SessionManager.SessionStart(c.Writer, c.Request)
			ss.Set("accountNum", u.StuNo)
			c.Redirect(http.StatusFound, "/userSpace")
		} else {
			fmt.Println(u.Name + " login wrong")
			c.HTML(http.StatusOK, "login.html", gin.H{
				"alert": tools.JsAlert("用户名或密码错误，请重新登录！"),
			})
		}
	}
}

func UsersUploadPost(c *gin.Context) {
	file, err := c.FormFile("usersFile")
	if err != nil {
		fmt.Println("a user's upload has err, please check what has happened.")
		return
	}
	dst := fmt.Sprintf("./Files/UsersUpload/%s", file.Filename)
	c.SaveUploadedFile(file, dst)
	ss, _ := SessionManager.CheckSession(c.Request)
	c.HTML(http.StatusOK, "upload.html", gin.H{
		"ifm":              tools.H5trans("<p><i>Upload success!</i></p>"),
		"userHeadImageDir": UASManager.GetUASFileDir(ss.Get("accountNum").(string), "HeadImage", db),
	})
}

func main() {

	defer db.Close()

	db.AutoMigrate(&model.SDUSTUser{})
	db.AutoMigrate(&model.UsersMessage{})
	db.AutoMigrate(&model.UserAddrSpace{})

	router := gin.Default()

	router.Static("staticFile", "./statics")
	router.Static("root", "./")
	router.StaticFS("/download", http.Dir("Files"))
	router.StaticFS("/downloadMovies", http.Dir("D:\\Movies"))

	router.LoadHTMLGlob("templates/*") //load all files; if has child dir use like "templates/**/*"

	router.GET("/indexPage", IndexPage)

	router.GET("/", Root)

	router.GET("/loginPage", LoginPage)

	router.GET("/uploadPage", tools.AuthMiddleWare(SessionManager), UploadPage)

	router.GET("/test", func(c *gin.Context) {
		c.HTML(http.StatusOK, "userMessageBox.html", gin.H{
			"model": "test1",
			"time":  "2020",
		})
	})

	router.GET("/t2", func(c *gin.Context) {
		c.HTML(http.StatusOK, "test.tmpl", nil)
	})

	// 已经设计了中间件，无需再次检查，只要取出session即可
	userSpace := router.Group("/userSpace", tools.AuthMiddleWare(SessionManager))
	{
		userSpace.Static("staticFile", "./statics") //路由组内重新设置静态路径

		// 用户主界面
		userSpace.GET("", UserSpaceRoot)

		// 用户消息界面
		userSpace.GET("/message", UserSpaceMessage)

		// 用户发送消息表单到服务器
		userSpace.POST("/message/sendMSG", UserSpaceMessageSendMSGPost)

		// 用户修改头像，上传文件到独立空间
		userSpace.POST("/message/headImage", tools.AuthMiddleWare(SessionManager), UserSpaceHeadImagePost)
	}

	router.GET("/downloadCenter", tools.AuthMiddleWare(SessionManager), DownLoadCenter)

	router.POST("/sign", SignPagePost)

	router.POST("/login", LoginPagePost)

	router.POST("/usersUpload", tools.AuthMiddleWare(SessionManager), UsersUploadPost)

	router.POST("/secret", func(c *gin.Context) {
		str := c.PostForm("psw")
		fmt.Println(tools.HashSecret(str))
	})

	router.Run(":9090")
}
