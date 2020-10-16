package main

import (
	tools "Gin/gorm2/gotools"
	model "Gin/gorm2/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"
)

// ding -config=ding.cfg -subdomain=zsc 8080
var (
	ServerIP string
	Version  string = "version_1.1.1"
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

// 主页/注册页
func IndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "newIndexPage.html", gin.H{
		"alert": "null",
	})
}

// 登录页/个人空间页，如若检查session存在则跳转个人空间页，否则进入登录页
func Root(c *gin.Context) {
	_, exist := SessionManager.CheckSession(c.Request)
	if exist {
		c.Redirect(http.StatusTemporaryRedirect, "/userSpace")
		return
	}
	c.Redirect(http.StatusFound, "/loginPage")
}

// 登陆页面
func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "newLoginPage.html", gin.H{
		"alert": "null",
	})
}

// 上传页面
func UploadPage(c *gin.Context) {
	ss, _ := SessionManager.CheckSession(c.Request)
	c.HTML(http.StatusOK, "newUploadPage.html", gin.H{
		"userHeadImageDir": UASManager.GetUASFileDir(ss.Get("accountNum").(string), "HeadImage", db),
	})
}

// 用户空间主页
func UserSpaceRoot(c *gin.Context) {
	ss, _ := SessionManager.CheckSession(c.Request)
	u := tools.CheckUserFromSession(ss, db)
	c.HTML(http.StatusOK, "newUserSpace.html", gin.H{
		"finish":           u.CheckUsersIFMFinish(),
		"userHeadImageDir": UASManager.GetUASFileDir(ss.Get("accountNum").(string), "HeadImage", db),
		"name":             u.Name,
		"stuNo":            u.StuNo,
		"domNo":            u.DomNo,
		"eMail":            u.Email,
		"age":              u.Age,
		"versionScript":    tools.GetVersionScript(Version),
	})
}

// 用户空间消息页
func UserSpaceMessage(c *gin.Context) {
	ss, _ := SessionManager.CheckSession(c.Request)
	u := tools.CheckUserFromSession(ss, db)
	fmt.Println(u)
	msg := make([]model.UsersMessage, 0, 0)
	db.Where("Receiver = ?", u.StuNo).Find(&msg)
	c.HTML(http.StatusOK, "newUserMessageBox.html", gin.H{ //转换为string类型以模板的形式传递给前端
		"userHeadImageDir": UASManager.GetUASFileDir(ss.Get("accountNum").(string), "HeadImage", db),
		"messages":         msg,
	})
}

// 用户发送消息
func UserSpaceMessageSendMSGPost(c *gin.Context) {
	ss, _ := SessionManager.CheckSession(c.Request)
	u := tools.CheckUserFromSession(ss, db)
	r := model.SDUSTUser{}
	db.Find(&r, "stu_no = ?", c.PostForm("MSG_receiver")) //从提交的表单中获取接收者信息
	alert := ""                                           //查找不到反馈
	if r.Name == "" {                                     // 无效的消息重置页面
		msgR := make([]model.UsersMessage, 0, 0)
		db.Where("Receiver = ?", u.StuNo).Find(&msgR) // 从数据库中查找当前用户的消息列表
		alert = "用户不存在"
		c.HTML(http.StatusOK, "newUserMessageBox.html", gin.H{
			"userHeadImageDir": UASManager.GetUASFileDir(ss.Get("accountNum").(string), "HeadImage", db),
			"messages":         msgR,
			"alert":            alert,
		})
	} else { // 有效则保存消息并重定向
		msgS := model.UsersMessage{
			Model:    gorm.Model{},
			Sender:   u.StuNo,
			Receiver: c.PostForm("MSG_receiver"),
			Content:  c.PostForm("MSG_content"),
			SendTime: time.Now(),
			READ:     false,
		}
		db.Create(&msgS)
		c.Redirect(http.StatusFound, "/userSpace/message")
	}
	//c.Redirect(http.StatusFound, "")
}

// 用户提交修改头像
func UserSpaceHeadImagePost(c *gin.Context) {
	file, err := c.FormFile("headImage")
	if err != nil {
		fmt.Println(tools.PrintLogHead(), "there is an err happened when upload the image.")
		return
	}
	// 获取session
	ss, _ := SessionManager.CheckSession(c.Request)
	// 根据session获取sid
	uasID := ss.Get("accountNum")

	// 应该交由UASManager来执行存储过程和数据库过程
	dst := UASManager.SaveUserFile(file, uasID.(string), db)

	// 按照指定位置保存并重定向
	c.SaveUploadedFile(file, dst)
	c.Redirect(http.StatusFound, "/userSpace")
}

func UserSpaceCloud(c *gin.Context) {
	ss, _ := SessionManager.CheckSession(c.Request)
	fp := UASManager.GetUASCloudPaths(ss.Get("accountNum").(string), db)
	c.HTML(http.StatusOK, "newUserCloud.html", gin.H{
		"userHeadImageDir": UASManager.GetUASFileDir(ss.Get("accountNum").(string), "HeadImage", db),
		"files":            fp,
	})
}

func UserSpaceCloudDownload(c *gin.Context) {
	dir := c.PostForm("dir")
	ss, _ := SessionManager.CheckSession(c.Request)
	file, err := os.Open("./" + UASManager.GetFileDir() + "/" + tools.CheckUASFromDB(ss.Get("accountNum").(string), db).UserAddr + "/" + dir)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fileName := path.Base(dir)
	fileName = url.QueryEscape(fileName)
	output, _ := ioutil.ReadAll(file)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	c.Data(http.StatusOK, "Content-Type", output)
}

// 下载中心页面
func DownLoadCenter(c *gin.Context) {
	ss, _ := SessionManager.CheckSession(c.Request)
	type dir struct {
		DirName  string
		RealAddr string
	}
	dirs := []dir{
		{"其他人上传的文件", "/download"},
		{"主站文件", "/downloadMovies"},
	}
	fmt.Println(dirs)
	c.HTML(http.StatusOK, "newDownloadPage.html", gin.H{
		"userHeadImageDir": UASManager.GetUASFileDir(ss.Get("accountNum").(string), "HeadImage", db),
		"dirCount":         2,
		"dirs":             dirs,
	})
}

// 注册提交
// 检查用户信息是否有重复
// 对用户密码加密
// 数据库存储用户信息表
// 初始化session
// 新建uas 数据库存储uas
func SignPagePost(c *gin.Context) {
	u := model.SDUSTUser{
		Model:    gorm.Model{},
		Name:     c.PostForm("UserName"),
		Password: tools.OnLock(c.PostForm("UserPwd"), c.PostForm("UserStuNo")),
		Age:      0,
		StuNo:    c.PostForm("UserStuNo"),
		Email:    c.PostForm("UserEmail"),
		DomNo:    c.PostForm("UserDomNo"),
		SID:      0,
	}
	u.Age, _ = strconv.Atoi(c.PostForm("UserAge"))
	if u.Name == "" || u.Password == "" {
		c.HTML(http.StatusOK, "newIndexPage.html", gin.H{
			"alert": "err",
		})
		return
	}
	db.Create(&u)
	ss := SessionManager.SessionStart(c.Writer, c.Request)
	ss.Set("accountName", u.Name)
	ss.Set("accountNum", u.StuNo)
	uas := UASManager.InitUserAddrSpace(u.StuNo)
	db.Create(&uas)
	c.Redirect(http.StatusFound, "/")
}

// 登录页提交
func LoginPagePost(c *gin.Context) {
	if c.PostForm("stuNo") == "" {
		c.HTML(http.StatusOK, "newLoginPage.html", gin.H{
			"alert": tools.JsAlert("学号不为空"),
		})
	} else {
		u := model.SDUSTUser{}
		db.Find(&u, "stu_no = ?", c.PostForm("stuNo"))
		fmt.Println(u.Password == tools.OnLock(c.PostForm("pwd"), c.PostForm("stuNo")))
		if u.Name != "" && u.Password == tools.OnLock(c.PostForm("pwd"), c.PostForm("stuNo")) { //登陆成功刷新Cookie/Session
			ss := SessionManager.SessionStart(c.Writer, c.Request)
			ss.Set("accountNum", u.StuNo)
			c.Redirect(http.StatusFound, "/userSpace")
		} else {
			fmt.Println(u.Name + " login wrong")
			c.HTML(http.StatusOK, "newLoginPage.html", gin.H{
				"alert": "err",
			})
		}
	}
}

// 用户上传提交
func UsersUploadPost(c *gin.Context) {
	file, err := c.FormFile("usersFile")
	resend := "success"
	if err != nil {
		fmt.Println("a user's upload has err, please check what has happened.")
		resend = "err"
	} else {
		dst := fmt.Sprintf("./Files/UsersUpload/%s", file.Filename)
		c.SaveUploadedFile(file, dst)
	}
	ss, _ := SessionManager.CheckSession(c.Request)
	c.HTML(http.StatusOK, "newUploadPage.html", gin.H{
		"ifm":              resend,
		"userHeadImageDir": UASManager.GetUASFileDir(ss.Get("accountNum").(string), "HeadImage", db),
	})
}

// 404
func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"msg": "地址不存在",
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

	// 404 response
	router.NoRoute(NotFound)

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
		userSpace.Static("root", "./")

		// 用户主界面
		userSpace.GET("", UserSpaceRoot)

		// 用户消息界面
		userSpace.GET("/message", UserSpaceMessage)

		// 用户发送消息表单到服务器
		userSpace.POST("/message-sendMSG", UserSpaceMessageSendMSGPost)

		// 用户修改头像，上传文件到独立空间
		userSpace.POST("/message-headImage", UserSpaceHeadImagePost)

		userSpace.GET("/cloud", UserSpaceCloud)

		userSpace.POST("/cloud-download", UserSpaceCloudDownload)
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
