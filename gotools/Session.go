package gotools

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

// Manager是一个Session管理类，创建时唯一指派对应的Cookie值，用于和SessionID匹配
// lock用于互斥访问
// provider是Session的底层实现方法簇，是一个接口，不同的实现方法可以生成不同的底层支持
// 最大存活时间用于判断Session是否过期以及进行GC
type Manager struct {
	cookieName  string     // private cookieName
	lock        sync.Mutex // protects session
	provider    Provider
	maxLifeTime int64
}

// 定义了一个接口，该接口用于根据服务器需要进行设计具体的底层实现
// 该接口是作为底层Session实现支持的方法簇，只规定了底层支持的行为，内容需要具体设计
type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionCheck(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
}

// MapProvider就是其中的一种Provider接口的实现
// 依托内存进行管理
// 抽象类型为map，sessionId是key， 对应唯一的session实例
type MapProvider struct {
	hash map[string]Session
}

// 初始化mp
func (mp *MapProvider) MPInit() {
	mp.hash = make(map[string]Session)
}

// 初始化一个Session
// 与上面的不同，上面的初始化是为了给其成员赋初值
// 而此处的初始化是接收到一个sessionId来初始化一个Session实例
func (mp *MapProvider) SessionInit(sid string) (ss Session, err error) {
	// 接收sid并新建一个Session对象，这里是mapProvider，底层实现是内存中的map 需要将这个Session实例插入到其中
	ss = &SimSession{sessionID: sid}
	ss = ss.(Session)
	mp.hash[sid] = ss
	return
}

// 检查Session 是否存在，不存在则返回一个新的Session以及错误，存在则直接返回
// 读取对应Id的Session并返回
// 由于会对不存在的情况创建新的Session导致该方法并不常用
// 如若不想进行新建的行为可以制作检查（下一方法）
func (mp *MapProvider) SessionRead(sid string) (ss Session, err error) {
	v, ok := mp.hash[sid]
	if !ok {
		ss, err = mp.SessionInit(sid)
		err = &SessionERR{errMSG: "Session NOT exist but we init a new one."}
		return
	}
	ss = v
	return
}

// 检查是否存在
// 仅返回结果，并不新建Session
func (mp *MapProvider) SessionCheck(sid string) (ss Session, err error) {
	v, ok := mp.hash[sid]
	if !ok {
		err = &SessionERR{errMSG: "Session NOT exist."}
	}
	ss = v
	return
}

// 移除目标Session
func (mp *MapProvider) SessionDestroy(sid string) error {
	delete(mp.hash, sid)
	return nil
}

// GC 并未完成
func (mp *MapProvider) SessionGC(maxLifeTime int64) {
	// todo
}

var Provides = make(map[string]Provider)

// 用于注册不同的底层实现
func Register(name string, provider Provider) {
	if provider == nil {
		panic("session: Register provider is nil")
	}
	if _, dup := Provides[name]; dup {
		panic("session: Register called twice for provider " + name)
	}
	Provides[name] = provider
}

// Session接口， 定义了Session的具体行为
type Session interface {
	Init()                            // init session's value
	Set(key, value interface{}) error // set session value
	Get(key interface{}) interface{}  // get session value
	Delete(key interface{}) error     // delete session value
	SessionID() string                // back current sessionID
}

// 实现了Session接口的一种 Session
type SimSession struct {
	// SimSession 索要维护的数据这里只有用户的账号，通过账户名来进行用户数据的维护
	sessionID string
	info      map[interface{}]interface{}
}

func (ss *SimSession) Init() {
	ss.info = make(map[interface{}]interface{})
}

func (ss *SimSession) Set(key, an interface{}) error {
	ss.info[key] = an
	return nil
}

func (ss *SimSession) Get(key interface{}) interface{} {
	if ss.info["accountNum"] == "" {
		return "User has not login."
	}
	return ss.info[key]
}

func (ss *SimSession) Delete(sid interface{}) error {
	// todo
	return nil
}

func (ss *SimSession) SessionID() string {
	return ss.sessionID
}

// provideName : 依托的session存储方式，用于获取对应的处理接口
func NewManager(provideName, cookieName string, maxLifeTime int64) (*Manager, error) {
	provider, ok := Provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
	}
	return &Manager{provider: provider, cookieName: cookieName, maxLifeTime: maxLifeTime}, nil
}

// 获取一个唯一的sessionId
func (manager *Manager) sessionId() string {
	b := make([]byte, 32)                   // 设置ID长度为32的字节数组
	if _, err := rand.Read(b); err != nil { // 赋随机值
		return ""
	}
	// 处理成编码
	return base64.URLEncoding.EncodeToString(b)
}

func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		sid := manager.sessionId()
		session, _ = manager.provider.SessionInit(sid)
		session.Init()
		cookie := http.Cookie{
			Name:     manager.cookieName,
			Value:    url.QueryEscape(sid),
			Path:     "/",
			HttpOnly: true,
			MaxAge:   int(manager.maxLifeTime),
		}
		http.SetCookie(w, &cookie)
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ = manager.provider.SessionRead(sid)
	}
	return
}

// 用manager中存有的Cookie键来取出上下文中存有的Cookie
// 这一步称之为取钥
// 再用Cookie存取的value来解码
// 用解码后的sid和manager中的provider匹配，检查是否存在这个sid
// 这一步称之为匹钥
// 只有匹钥成功了才能确定该用户是有效用户
func (manager *Manager) CheckSession(r *http.Request) (ss Session, exist bool) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	}
	sid, _ := url.QueryUnescape(cookie.Value)
	if s, err := manager.provider.SessionCheck(sid); err != nil {
		return
	} else {
		ss = s.(Session)
		exist = true
	}
	return
}

// 实现了一个error
type SessionERR struct {
	errMSG string
}

func (se *SessionERR) Error() string {
	return se.errMSG
}
