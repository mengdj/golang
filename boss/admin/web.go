package admin

import (
	"context"
	"encoding/json"
	"ext"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/alecthomas/log4go"
	"github.com/gin-gonic/gin"
	"github.com/go-session/gin-session"
	"github.com/gogf/gf/container/gtype"
	"github.com/panjf2000/ants/v2"
	"gopkg.in/olahol/melody.v1"
	"net/http"
	"tool"
	"admin/model"
)

//全局管理员
var globalUser = &model.User{Account: tool.DEFAULT_ADMIN, Password: ""}
var serverPort string
var serverAddr string

type Status struct {
	Code int8        `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type Ping struct {
	Query    string      `json:"query"`
	Response interface{} `json:"response"`
}

type Web struct {
	*gin.Engine
	logger    *log4go.Logger
	mill      *ext.ExtGoChanel
	goroutine *ants.Pool
	bcount    *gtype.Int32
}

func NewWeb(logger *log4go.Logger, m *ext.ExtGoChanel, p *ants.Pool) *Web {
	return &Web{Engine: gin.New(), logger: logger, mill: m, goroutine: p, bcount: gtype.NewInt32(0)}
}

func (this *Web) Run(addr ...string) (err error) {
	var (
		msgOnline, msgOffline <-chan *message.Message
		lody                  *melody.Melody
	)
	//基本配置
	lody = melody.New()
	gin.SetMode(gin.ReleaseMode)
	this.Use(gin.Recovery())
	this.Use(gin.Logger())
	this.Use(ginsession.New())
	if err = this.init(lody); nil != err {
		return err
	}
	globalUser.Password = tool.RandomStr(tool.RANDOM_CHARS, 12)
	if serverAddr, err = tool.LocalAddress(); nil != err {
		serverAddr = "127.0.0.1"
	}
	serverPort = addr[0]
	this.logger.Info("管理地址:http://%s%s", serverAddr, addr[0])
	this.logger.Info("管理账号:%s", globalUser.Account)
	this.logger.Info("管理密码:%s", globalUser.Password)
	//订阅消息(并处理消息)
	if msgOnline, err = this.mill.Subscribe(context.Background(), tool.WORKER_ONLINE_RESULT); nil != err {
		return err
	}
	if msgOffline, err = this.mill.Subscribe(context.Background(), tool.WORKER_OFFLINE_RESULT); nil != err {
		return err
	}
	//先暂停此消息队列的发送，待管理员上线唤醒（节约资源）
	this.mill.Pause(tool.WORKER_ONLINE, tool.WORKER_OFFLINE)
	//广播WORKER上线消息(任何一个客户端上线都会触发此操作)
	this.goroutine.Submit(func() {
		for msg := range msgOnline {
			worker:=model.Worker{}
			if nil == json.Unmarshal(msg.Payload, &worker) {
				if res, err := json.Marshal(Ping{Query: tool.WORKER_ONLINE_RESULT, Response: worker}); nil == err {
					lody.Broadcast(res)
				}
			}
			msg.Ack()
		}
	})
	//广播WORKER离线消息(任何一个客户端离线都会触发此操作)
	this.goroutine.Submit(func() {
		for msg := range msgOffline {
			worker:=model.Worker{}
			if nil == json.Unmarshal(msg.Payload, &worker) {
				if res, err := json.Marshal(Ping{Query: tool.WORKER_OFFLINE_RESULT, Response: worker}); nil == err {
					lody.Broadcast(res)
				}
			}
			msg.Ack()
		}
	})
	return this.Engine.Run(addr...)
}

func (this *Web) init(lody *melody.Melody) (err error) {
	this.Delims("{{", "}}")
	this.Static("/static", "./web/view/static")
	this.Static("/data", "./data")
	this.StaticFile("/favicon.ico", "./web/view/favicon.ico")
	this.LoadHTMLFiles("./web/view/index.html", "./web/view/login.html")
	//处理websocket()
	lody.HandleConnect(func(session *melody.Session) {
		//开始连接
		this.bcount.Add(1)
		this.mill.Continue(tool.WORKER_ONLINE, tool.WORKER_OFFLINE)
	})
	lody.HandleDisconnect(func(session *melody.Session) {
		//丢失连接
		this.bcount.Add(-1)
		if this.bcount.Val() == 0 {
			//暂时屏蔽信息
			this.mill.Pause(tool.WORKER_ONLINE, tool.WORKER_OFFLINE)
		}
	})
	lody.HandleMessage(func(session *melody.Session, bytes []byte) {
		//处理消息
		query:=string(bytes)
		switch query {
		case tool.QUERY_WORKERS:
			//查询所有在线客户端数据
			if !this.mill.IsPause(tool.QUERY_WORKERS) {
				this.goroutine.Submit(func() {
					if messages, err := this.mill.Subscribe(context.Background(), tool.QUERY_WORKERS_RESULT); nil == err {
						msg := <-messages
						workers := []model.Worker{}
						if nil == json.Unmarshal(msg.Payload, &workers) {
							if res, err := json.Marshal(Ping{Query: tool.QUERY_WORKERS_RESULT, Response: workers}); nil == err {
								session.Write(res)
							}
						}
						msg.Ack()
					}
				})
				this.mill.Publish(tool.QUERY_WORKERS, message.NewMessage(watermill.NewUUID(), []byte{}))
			}
		}
	})
	this.Use(crosHandler)
	this.GET("/", login)
	this.POST("/sign", sign)
	this.GET("/ping", func(c *gin.Context) {
		lody.HandleRequest(c.Writer, c.Request)
	})
	v1 := this.Group("/v1")
	//需要授权才能访问的
	v1.Use(authHandler)
	{
		v1.GET("/", index)
	}
	return
}

//登录数据
func login(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", gin.H{
		"username": tool.DEFAULT_ADMIN,
	})
}

//签名数据(完成JWT)
func sign(c *gin.Context) {
	var (
		usr    model.User
		status Status
		err    error
		token  string
	)
	err = c.ShouldBind(&usr)
	status.Code = -1
	if nil != err {
		status.Msg = err.Error()
		c.JSON(http.StatusBadRequest, status)
	} else if usr.Account != globalUser.Account || usr.Password != globalUser.Password {
		status.Code = -2
		status.Msg = "账号或密码错误"
		c.AsciiJSON(http.StatusOK, status)
	} else {
		if token, err = model.GetToken(usr.Account, usr.Password); nil == err {
			c.Header("token", token)
			if "" != c.Query("ajax") {
				status.Code = 0
				status.Msg = "登录成功"
				status.Data = token
				c.AsciiJSON(http.StatusOK, status)
			} else {
				//直接到登录页
				c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/v1?token=%s", token))
			}
		} else {
			status.Code = -3
			status.Msg = err.Error()
			c.AsciiJSON(http.StatusServiceUnavailable, status)
		}
	}
}

//首页
func index(c *gin.Context) {
	if token, exist := c.Get("token"); exist {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"username": tool.DEFAULT_ADMIN,
			"client":   c.ClientIP(),
			"token":    token,
			"server":   serverAddr,
			"port":     serverPort,
		})
	} else {
		c.Redirect(http.StatusUnauthorized, "/")
	}
}

//授权检查中间件
func authHandler(c *gin.Context) {
	var token string
	if token = c.Query("token"); "" != token {
		usr, err := model.ParseToken(token)
		if nil != err {
			c.JSON(500, err)
			c.Abort()
		}
		if usr.Account != globalUser.Account || usr.Password != globalUser.Password {
			c.JSON(http.StatusUnauthorized, err)
			c.Abort()
		}
	} else {
		c.JSON(http.StatusUnauthorized, "未授权的操作，请重新登录")
		c.Abort()
	}
	c.Set("token", token)
	c.Next()
}

//跨域访问：cross  origin resource share
func crosHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Origin", "*") // 设置允许访问所有域
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
	c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma,token,openid,opentoken")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
	c.Header("Access-Control-Max-Age", "172800")
	c.Header("Access-Control-Allow-Credentials", "false")
	c.Next()
}
