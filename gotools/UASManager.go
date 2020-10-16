package gotools

import (
	model "Gin/gorm2/models"
	"fmt"
	"github.com/jinzhu/gorm"
	"mime/multipart"
	"sync"
)

// UAS 管理类，用于管理具体的UAS实例
type UASManager struct {
	lock         sync.Mutex
	fileDir      string
	addrPreFix   string
	MAXSpaceSize int64
}

// 新建一个UAS管理类的方法， 第一个参数是当前类型UAS的根目录， 第二个参数是UAS地址空间前缀， 最后一个参数是当前类型的UAS空间最大值
func NewUasManager(fd string, addrPF string, maxSS int64) *UASManager {
	return &UASManager{
		lock:         sync.Mutex{},
		fileDir:      fd,
		addrPreFix:   addrPF,
		MAXSpaceSize: maxSS,
	}
}

// 初始化一个UAS，依据是唯一的uid，返回一个实例
func (um *UASManager) InitUserAddrSpace(uid string) *model.UserAddrSpace {
	um.lock.Lock()
	defer um.lock.Unlock()
	// 1.生成一个新的用户地址空间
	// 2.根据manager的值初始化对应的用户地址空间变量
	// 3.具体申请硬盘空间，建立文件区，初始化文件区配置
	// 4.返回对象交给上层创建数据库对应记录
	// *.如果可以将添加回滚效果
	newUAS := model.UserAddrSpace{}
	newUAS.Set(uid, um.CreateUASAddrDir(uid), 0, um.MAXSpaceSize)

	// 申请建立文件区
	CreateDir(um.fileDir+"/"+newUAS.UserAddr, "root")
	CreateDir(um.fileDir+"/"+newUAS.UserAddr, "headImg")
	// TODO 其他文件夹建立，配置任务

	return &newUAS
}

// 创建用户地址空间根目录（目录名）参数是唯一的uid
func (um *UASManager) CreateUASAddrDir(uid string) string {
	return um.addrPreFix + uid
}

// 获取当前管理器的所有UAS的父目录， 也就是管理器管理的根目录
func (um *UASManager) GetFileDir() string {
	return um.fileDir
}

// 获取一个UAS的指定子目录
func (um *UASManager) GetUASFileDir(uid string, fileDirName string, db *gorm.DB) string {
	fd := um.fileDir
	uas := CheckUASFromDB(uid, db)
	fd = "root/" + fd + "/" + uas.UserAddr + "/" + fileDirName
	return fd
}

// 获取当前用户的云空间目录 调用save包的功能 返回一个切片
func (um *UASManager) GetUASCloudPaths(uid string, db *gorm.DB) []FilePath {
	fd := um.GetUASFileDir(uid, "", db)
	pfL := len(um.fileDir) + len(CheckUASFromDB(uid, db).UserAddr)
	return WalkThroughDir(fd[5:], pfL)
}

// 保存指定用户的某个文件 返回一个保存路径，同时在数据库中刷新数据内容
func (um *UASManager) SaveUserFile(file *multipart.FileHeader, uid string, db *gorm.DB) (dst string) {
	// 从数据库查找uid对应的用户数据记录
	uas := &model.UserAddrSpace{}
	FindUAS(uas, "user_id = ?", uid, db)
	// 创建路径目录
	dst = fmt.Sprintf("%s/%s/HeadImage/%s", um.GetFileDir(), uas.UserAddr, "headImage.jpg")

	// 获得用户空间目录大小
	uas.CurrentSpace, _ = GetDirSize(um.GetFileDir() + "/" + uas.UserAddr)
	if uas.CurrentSpace < 14 {
		uas.CurrentSpace += file.Size
	}
	// 从数据库更新当前值
	UpdateUASSingle(uas, "current_space", uas.CurrentSpace, db)

	fmt.Println("Now users current space is ", uas.CurrentSpace, " / ", uas.MAXSpace)
	return
}
