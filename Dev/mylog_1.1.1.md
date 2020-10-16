> 2020年10月8日

`[MOD]`修改了*注册页面*的**项目条款**部分，将原来的特效修改为*modal*模态框效果，依旧是用`Javascript`来修改其内部数据(innerHTML)。  

------  
> 2020年10月9日

`[NEW]`添加了*登陆页面*的账号/密码验证控件，采用了bootStrapValidator控件，经**离线/在线**测试功能完整安全，可以满足当前使用。  
`[MOD]`修改了*登陆页面*的服务器返回错误信息方法，将原来的服务器返回JavaScript代码方式修改为：仅返回错误代码，由前端根据代码监听并调用生成模态框。减小了服务返回数据的量，同时规范了前端行为规范。  
`[MOD]`修改了*登陆页面*和*注册页面*的导航栏中的‘主页’按钮，将其转到地址修改为<kbd>href = '/indexPage'</kbd>。  
`[MOD]`修改了*用户主页*的日志模块排版细节。  
`[OVR]`重写了*用户主页*的JavaScript内容设计方式，将原来的直接调用模板对象成员修改为：‘统一由JavaScript代码接收模板对象，存储为一个变量，在使用到的时候访问其字段获取值’的方式。减小了耦合性，统一了前端行为规范，将此设计方法推广到前端的所有页面设计当中，作为项目的规范方法。
  
------  
> 2020年10月11日  

`[MOD]`修改了*注册页面*的服务器返回错误信息方法，将原来的服务器返回JavaScript代码方式修改为：仅返回错误代码，由前端根据代码监听并调用生成模态框。减小了服务器返回数据量，规范了前端行为规范。将此规范作为新的前后端数据交互规范。  
`[STD#01]`项目开发规范标准#01:  
服务器提交给前端数据仍旧以模板的方式交互，在后端代码中规范返回的数据内容，以不同的字段代表数据，其数据集合为一个标准的对象；前端使用代码<kbd>var msg = {{ . }}</kbd>获取总的对象实例，根据后端定义的数据内容获取指定字段的数据，对前端的功能、属性进行对应的部署。通过这种规范标准化前后端的工作内容，简化工作步骤，减小两者之间的交互复杂性。代表性的应用是模态框接收alert警告的处理方式（见注册页和登录页）。  

------
> 2020年10月12日

`[OVR]`重写了*注册页面*的JavaScript设计方式，采用了**STD#01**标准。  
`[OVR]`重写了*登录页面*的JavaScript设计方式，采用了**STD#01**标准。  
`[OVR]`重写了*下载中心*的Go服务器代码，采用了**STD#01**标准。  
新的设计在模板数据对象里有三个字段信息：头像地址、需要解析的下载地址个数、地址列表。所有字段无需序列化即可直接传递使用，修改时只要更改数据项的地址和数据项个数即可。前端获取模板对象后根据具体字段的名即可获取到对应的数据值。前端的表格生成部分也由固定长度改为根据数据动态生成，更为方便。  
```go
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
```
```javascript
var msg = {{.}}
document.getElementById("hi").src = msg['userHeadImageDir'] + '/headImage.jpg'
console.log(msg)
var dirs = msg['dirs']
var dirCount = msg['dirCount']
```  
`[OVR]`重写了*上传页面*的JavaScript代码设计方式，采用了**STD#01**标准。
`[MOD]`修改了*上传页面*的Go服务器代码中的模板数据格式，取消原来的返回HTML代码的方式，改为返回状态码，由前端根据状态码执行操作。服务器代码根据文件上传成功与否修改状态码，最后统一返回数据。  
`[OVR]`重写了*用户消息页面*的Go服务器代码，修改了模板返回对象的数据格式，从原来的内嵌json格式调整为全部为模板数据格式，读取的消息数据列表直接放在返回的模板数据中，前端可以直接调用。  
字段名定义如下：
```go
    msg := make([]model.UsersMessage, 0, 0)
    c.HTML(http.StatusOK, "newUserMessageBox.html", gin.H{
        "userHeadImageDir": "/xx/xx",
        "messages":         msgR,
        "alert":            alert,
    })
```  
`[OVR]`重写了*用户消息页面*的JavaScript代码设计方式，根据最新的后端返回模板数据格式采用**STD#01**标准重新编写。追加了对于没有数据时的信息反馈。  
`[MOD]`修改了*注册页面*和*登陆页面*的模板数据格式，给基本的GET请求添加了一个<kbd>"alert": "null"</kbd>内容，解决了验证控件不成功的问题。  
`[NEW]`追加了*注册页面*的信息验证，对每一种信息在前端进行数据验证格式性确认，而重复性仍旧在后端完成。  
`[RMV]`删除了*注册页面*的项目条款check控件，取而代之的是会员特权说明。鉴于项目条款对于本项目没有太大帮助所以移除该功能。  
`[MOD]`修改了`/`目录的处理函数，将没有登录的用户将重定向到登录页，已经有Session的用户将回到个人主页。  

------  
> 2020年10月13日  

`[NEW]`开始设计*云空间*模块，未开发完成。Go服务器设计上，单独开了一个处理页面请求，当登陆过的用户访问该url时，UASManager会根据其uid获得到UAS空间地址。在`Save`包中设计了一个新的结构体，用于存储目录结构，定义如下：  
```go
type FilePath struct {
	FileType byte // 表示文件种类默认0为目录，其他为文件
	FileAddr string // 目录或文件名
}
```
Save包内定义的函数可以访问指定目录，并将沿途的目录内容记录下来，最后返回给调用处。  
```go
// 遍历root目录下所有文件，将结果放在一个切片里返回
func WalkThroughDir(root string) []FilePath {
	fmt.Println(root)
	fp := make([]FilePath, 0)
	walk := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			fp = append(fp, FilePath{
				FileType: 0,
				FileAddr: path,
			})
		} else {
			fp = append(fp, FilePath{
				FileType: 1,
				FileAddr: path,
			})
		}
		return nil
	}

	err := filepath.Walk("./"+root, walk)
	if err != nil {
		log.Fatal(err)
	}
	return fp
}
```
在调用生成目录后，将结构体发送给前端，前端获得对应的目录，动态生成地址，可以是a标签，以提供下载服务。这个模块还需要继续设计。  
除了这部分内容，针对云空间的下载服务还设计了指定目录文件的下载服务，同样被设计在一个页面上。设计了一处*POST*请求，用以接受一个目录地址，只可以是文件地址，设计的*handlefunc*为：`/downloadTest(暂时)`，提取出post请求的地址，读取文件并修改*response*的头部信息，将文件转换为字节流放在*response*中返回。以这种方式来提供特定文件的下载服务。该模块也未设计完成。  
`[MOD]`修改*用户主页*的头向下拉按钮中的‘云空间’转到地址为`'/dlt'(暂时)`
  
-----

> 2020年10月14日 --开始  
2020年10月15日 --继续完善

`[RCH]`完善了所有登陆后的*权限页面*的导航条中‘我的’一项的下拉列表内容，追加了‘云空间’的地址，增加了新的图标。  
`[NEW]`编写前端目录结构框架。  
### 数据结构定义
定义节点结构，用于表示一个目录系统中的节点，该节点可以是文件，也可以是目录（目录作为特殊的文件）；节点拥有三个字段：名字、类型、子目录（可选）。名字是唯一标识该节点的字段，采用全地址作为名字。类型标识该节点的种类，如果是目录节点该值为0，其他类型为非0。种类为目录的节点可以拥有子目录，其余种类为`null`  
```javascript
/* 新建节点
** 节点类型定义：-1为根节点 0为目录节点 1为文件节点
** 根节点只有一个，创建时直接在字典中注册
** 目录节点会生成子目录
** 文件节点不生成子目录
** 非根节点会找寻父节点并挂在父节点下
** 最后注册当前节点
*/
function TreeNode(name, type){
	this.Name = (name === undefined ? null : name)
	this.Type = (type === undefined ? 1 : type)
	if(this.Type == -1){
		dict['root'] = this
		this.Child = []
		return
	}else if(this.Type == 0)
		this.Child = []
	else
		this.Child = null
	SetTreeFather(this)
	if(this.Type == 0)
		dict[this.Name] = this
}

/* 设置节点的父节点，查询最后一个斜杠的位置并分析前部分切片的长度
** 若是根目录子节点则添加节点到根目录下
** 其他节点挂在对应注册点的节点子目录下
*/
function SetTreeFather(node){
	var name = node.Name
	for(var i=name.length; i>=0 && name[i]!='\\'; i--);
	if(i >= 0){
		var str = name.slice(0, i)
		dict[str].Child.push(node)
	}else{
		dict['root'].Child.push(node)
	}
}
```
### 行为定义  
#### 生成树结构
前端将会接收到一个严格递归顺序的目录序列，按照该顺序生成每个目录中的节点，并在dict中注册。  
前端将收到一个顺序列表，每个元素的类型定义如下：  	
```go
type FilePath struct {
	FileType byte   // 表示文件种类默认0为目录，其他为文件
	FileAddr string // 目录或文件名
}
```
由go访问服务器目录获得的一个目录数组，传递给前端的js处理，js收到的是一个对象数组，由于顺序是严格递归定义的，所以就可以逐个遍历生成节点，创建节点之间的关系，最终就可以获得一整个目录树结构。  
注意：  
第一个节点是根节点，需要特殊处理，其FileType字段强制修改为-1。
#### 渲染
对于每个目录下渲染当前目录结构到HTML中，每当发生目录跳转时会渲染新的目录结构。  
当所有节点正确生成后即可开始渲染，以`root`节点进行初步渲染。  
渲染的对象为某一个目录节点，表格的第一行为其父节点的项目。其余行数据为其`Child`字段的所有数据，针对每个数据元素的内容生成节点类型、节点名/路径。  
每个数据单元的类型为其节点的类型，而其后的路径就需要针对是目录或者文件来异化。  
文件：生成a标签，onClick行为是点击获取其值，输出到某个位置。  
目录：生成a标签，onClick行为是更改目录，发生目录切换。
#### 目录切换
每当目录发生切换，包括回退到上级目录、跳转到子级目录等行为时会获取一个新的目录列表，此时执行渲染操作在表格中生成新的目录结构。  
### 设计成果
> 1.0版本  

制作出了简单的文件管理器，可以提供给用户查看自己的云空间目录效果，可以递归访问各个目录，查看文件、选择文件时可以将文件的地址复制到下载输入框中。将功能移植到主要文件中。  
由于前端获得的是屏蔽上层目录的文件结构，所以如果发送文件地址给后端是下载不到的，需要在后端处理时解析，增加前缀以访问到具体的文件，这也是保证用户身份安全的一个方法。  
`[NEW]`新建了*云空间*页面的Go代码，设计了接收下载地址的POST请求方法，该方法接收来自用户的url，通过增加前缀来访问真实的目录地址，调用os包来找到目标位置的文件，通过添加到Response的头部来将结果返回给客户端。  
`[NEW]`新设计了*UASManager*包的`GetUASCloudPaths`方法，该方法会调用*Save*包的`WalkThroughDir`函数来递归访问指定uid用户的目录，并将结果以一个*严格递归顺序* 的切片返回给上层。  
`[NEW]`新设计了*Save*包的*WalkThroughDir*函数，该函数将使用 *filepath.Walk()* 函数来递归访问指定目录，并以指定长度`pfL`来切割每个结点的路径结果，存储在一个切片中，最终将结果返回给调用者。由于访问是按照`dfs`的方式递归访问，所以返回的结果也是按照*严格递归顺序*排布的。  
`[MOD]`调整了*云空间*页面的GET、POST请求的handleFunc排布方式，将他们从主程序中脱离出来，写成函数的形式，减小耦合度。调用他们的位置修改为*UserSpace*路由组，对应的访问地址修改为`/cloud` 和 `/cloud-download`。  
`[MOD]`调整了*前端页面*的有关*云空间*的访问地址，将测试用的地址转换为现行地址。  