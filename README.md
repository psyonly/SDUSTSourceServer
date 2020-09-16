# SDUSTSourceServer
## 前言
鉴于学校缺少一个资源站（或者我不知道）决定设计一个校园资源服务器来给全校师生提供一个资源下载和共享的平台。  
编程语言我选择了*Go*语言，因为这一阶段都在学习Go语言，且得益于Go本身的Web包和各种开源框架能够较为简单的架设起整个环境。  
### 使用了哪些工具  
1.Golang 1.14  
2.Gin框架1.6.3  
3.gorm框架  
4.mysql  
### 功能介绍  
**main.go** 是主程序的入口，包括了服务器数据交互，GET/POST请求处理等多个基本功能。  
gotools 包内是包含了对Html/Javascript 进行转义或相应处理的方法。我在后续的更新中会持续为这个包添加新的方法。  
#### 版本信息  
v1.0.1在宿舍和机房调试后暂时都可以正常运行，多用户访问内网IP均可正常服务。  
*已有的功能*  
用户注册信息并保存至服务器。  
登录后可以在个人中心页面查看当前登录状态、资料完整度。  
一次登录即可在客户端保存Cookie[120s]。  
Cookie存活期内访问权限页面时服务器会根据Cookie的内容查询数据库有无对应账户，合法用户会允许访问，非法用户会禁止访问并返回JSON。（完成Cookie-Session双重认证功能）  
访问主站地址为IP+端口，如192.168.135.103:9090/ 或 192.168.135.103:9090。若登陆过并且Cookie活动时间还存在的用户再次访问该地址将自动重定向到个人中心页面。  
登录用户可以访问上传页面和下载页面，当前设置并无上传和下载次数限制。  

----
**v1.0.2**  
本次更新的功能：增设了获取服务器本地IP的功能，可以根据不的网络自动配置Cookie中的domain信息，不必每次手动输入IP。  
将获取IP的功能实现设计为函数，后续决定要不要单独设置一个包。  
修正了第一次注册登陆的用户无法存储Cookie的问题，目前的Cookie自动刷新只有老用户登录、新用户注册两个接口可以调用。  
setCookie 功能封装为一个方法，后续版本决定是否设置单独包和包合并的问题 **见H1**  
消息系统开发：依托于Go模板引擎和JS解析JSON，从服务器发送系列化后的消息数据数组（JSON），通过模板传递给前端，前端JS代码调用模板获取JSON字符串，  
将其用JS解析成可以识别的JS对象，呈现在页面上。  
用户点击主页的“我的消息”会跳转到一个链接，链接中以表格的形式显示用户的消息列表，目前设置四个字段：Sender Receiver Content Read。  
页面可以向任意用户发送消息:  
通过本页面的form输入接收者学号、信件内容即可发送消息，发送失败（无效的接收者）会重载页面并标记错误，成功则会刷新页面。  
>_H1:已创建Cookies包和SDUSTUser包用于存储跟Cookie相关的功能以及用户结构体及其相关方法  
------

**v1.0.3**  

计划加入的特性：密码加密功能，对接受的来自用户的密码进行加密算法，初步设想是自己设计，为tools包下的secret.go中，看看效果如何再移植，如果自己的算法效果不好就用已有的算法。  

重写了Cookie/Session部分，采用了更加安全的验证方法，通过设计一组方法簇来抽象出针对不同实现形式的Session管理器，依托管理器来对Session进行管理，完成Session的初始化、管理等；同时使Session不直接与数据库/底层存储交互，利用管理器来减少耦合性。  [TIPS]:Session开发不完备，对于GC和生命周期未完善。 

[REMAKE:移除Cookie包内容，将Cookie包内的部分功能和Session包合并，将原来的代码保存，修改所有API，设计新包DBCheck，用于实现各种管理器和数据库交互的功能。  

计划加入用户"独立文件空间"概念，并第一个提供用户头像功能。  
*独立地址空间设计:**UAS***  
给每个用户提供独立的文件夹目录用于存储用户的文件等信息，计划是如网盘软件一样限制每个用户的空间大小，给予一定分配。

所要做的任务：

- 用户第一次注册时创建好一个文件夹，用于放该用户的文件信息。
- 这部分功能需要依托OS来完成mkdir等工作.  指定文件夹的名称为唯一地址，并将地址存储于数据库的表中。
- 由于采用orm连接数据库表的方式，所以这里打算考虑结构的类型，直接利用表迁移功能确定表的结构  

```go
type UserAddrSpace struct{
	Model
	userID string
	userAddr string
	currentSpace uint64
	MAXSpace uint64
}
```

其中userID是可以区分一个用户的唯一标识，userAddr是分配给每个用户的根目录，对用户屏蔽该目录，只暴露子目录给用户，该根目录存储一些用户的文件，用来方便系统对其进行管理。  
最后两个字段用于记录当前使用程度。  

---

已设计完成的计划：  
添加UAS及其管理器部分功能，整合了model包以及更新了gotools包。  
添加了一个简单的日志生成工具，需要继续增加细节和整理。  
修改main函数中的handlFunc独立出来，但未完全独立，仍与main函数在同一个文件下，因为handleFunc中有一些管理器等参数，暂时无法直接剥离，会在后续开发中完善。  
凭借UAS实现了一个简单的用户头像功能，初次注册即可初始化UAS空间，并内置一个headImage文件夹（用户屏蔽此层）用于存放头像。  
其他部分作为用户个人空间。  

> 其他:前端部分已有新的开发人员帮助开发，期待更美观的界面。

------

**v1.0.4**

更新：修改了用户密码的存储方式，使用了加密算法使得用户数据更为安全。  
设计更新了`secret`包，在其中做安全存储的功能。其中有加密方法、验证方法和加盐常量设置。将此包重命名为`Secret`  

设计了`Save`包，内部存放一些关于对于文件的操作，通过设置此包来减少主程序与文件的直接操作。  

##### 其他
资源部份有五张Stardewvalley资源图片，若存在侵权请联系我删除。  
感谢Q1mi老师的go教程，受益匪浅。其博客地址liwenzhou.com  
