package admin

import (
	"admin/model"
	"context"
	"ext"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/alecthomas/log4go"
	"github.com/gin-gonic/gin"
	"github.com/go-session/gin-session"
	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/text/gstr"
	"github.com/json-iterator/go"
	"github.com/panjf2000/ants/v2"
	"gopkg.in/olahol/melody.v1"
	"net/http"
	"time"
	"tool"
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
	ctx       context.Context
}

func NewWeb(c context.Context, logger *log4go.Logger, m *ext.ExtGoChanel, p *ants.Pool) *Web {
	return &Web{Engine: gin.New(), ctx: c, logger: logger, mill: m, goroutine: p, bcount: gtype.NewInt32(0)}
}

func (this *Web) Run(addr ...string) (err error) {
	var (
		messages <-chan *message.Message
		lody     *melody.Melody
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
	this.logger.Info("结束服务请按组合键:CTRL+C")
	//订阅消息(并处理消息 WORKER_ADMIN_CTX_RESULT为反馈给web端的客户上下线通知)
	if messages, err = this.mill.Subscribe(this.ctx, tool.WORKER_ADMIN_CTX_RESULT); nil != err {
		return err
	}
	//先暂停此消息队列的发送，待管理员上线唤醒（节约资源）
	this.mill.Pause(tool.WORKER_ADMIN_CTX)
	//广播WORKER离线消息(任何一个客户端离线都会触发此操作)
	this.goroutine.Submit(func() {
		for {
			select {
			case msg := <-messages:
				msg.Ack()
				meta := msg.Metadata.Get(tool.WORKER_ADMIN_CTX_SUB_TYPE)
				switch meta {
				case tool.WORKER_ONLINE, tool.WORKER_OFFLINE:
					//上线 离线
					worker := model.Worker{}
					if tool.WORKER_ONLINE == meta {
						if nil == jsoniter.Unmarshal(msg.Payload, &worker) {
							if res, err := jsoniter.Marshal(Ping{Query: tool.WORKER_ONLINE, Response: worker}); nil == err {
								lody.Broadcast(res)
							}
						}
					} else {
						if nil == jsoniter.Unmarshal(msg.Payload, &worker) {
							if res, err := jsoniter.Marshal(Ping{Query: tool.WORKER_OFFLINE, Response: worker}); nil == err {
								lody.Broadcast(res)
							}
						}
					}
					break
				case tool.WORKER_QPS:
					//QPS
					qps := model.Qps{}
					if nil == jsoniter.Unmarshal(msg.Payload, &qps) {
						if res, err := jsoniter.Marshal(Ping{Query: tool.WORKER_QPS, Response: qps}); nil == err {
							lody.Broadcast(res)
						}
					}
					break
				}
				break
			case <-this.ctx.Done():
				//term
				goto EXIT_CUR
			}
		}
	EXIT_CUR:
		//code here
	})
	return this.Engine.Run(addr...)
}

func (this *Web) init(lody *melody.Melody) (err error) {
	this.Delims("{{", "}}")
	this.Static("/static", "./admin/view/static")
	this.Static("/data", "./data")
	this.StaticFile("/favicon.ico", "./admin/view/favicon.ico")
	this.LoadHTMLFiles("./admin/view/index.html", "./admin/view/login.html")
	//处理websocket()
	lody.HandleConnect(func(session *melody.Session) {
		//开始连接
		this.bcount.Add(1)
		this.mill.Continue(tool.WORKER_ADMIN_CTX)
		//创建独自的消息队列处理请求消息
	})
	lody.HandleDisconnect(func(session *melody.Session) {
		//丢失连接
		this.bcount.Add(-1)
		if this.bcount.Val() == 0 {
			//暂时屏蔽信息(所有管理都下线了就不再发送消息了)
			this.mill.Pause(tool.WORKER_ADMIN_CTX)
		}
	})
	lody.HandleMessage(func(session *melody.Session, bytes []byte) {
		//处理消息
		query := string(bytes)
		switch query {
		case tool.QUERY_WORKERS:
			//5秒钟没有返回结果则退出协程
			ctxTimeout, _ := context.WithTimeout(this.ctx, time.Second*5)
			this.goroutine.Submit(func() {
				if messages, err := this.mill.Subscribe(ctxTimeout, tool.WORKER_ADMIN_CTX_RESULT); nil == err {
					select {
					case msg := <-messages:
						msg.Ack()
						switch msg.Metadata.Get(tool.WORKER_ADMIN_CTX_SUB_TYPE) {
						case tool.QUERY_WORKERS:
							workers := []model.Worker{}
							if nil == jsoniter.Unmarshal(msg.Payload, &workers) {
								if res, err := jsoniter.Marshal(Ping{Query: tool.QUERY_WORKERS, Response: workers}); nil == err {
									//请求需要时间
									if !session.IsClosed() {
										session.Write(res)
									}
								}
							}
							break
						}

						break
					case <-ctxTimeout.Done():
						break
					}
				}
			})
			this.mill.Publish(tool.WORKER_ADMIN_CTX, ext.NewAdminSubTypeMessage([]byte{}, tool.QUERY_WORKERS))
			break
		case tool.QUERY_PING:
			//回复心跳消息
			if res, err := jsoniter.Marshal(Ping{Query: tool.QUERY_PING, Response: []byte("PONG")}); nil == err {
				session.Write(res)
			}
			break
		default:
			//处理其他信息
			if pos := gstr.Pos(query, tool.REFRESH_SELECTED_WORKER); -1 != pos {
				//刷新客户信息(REFRESH_SELECTED_WORKER:ip1;ip2;ip3)
				if ips := gstr.SubStr(query, pos+len(tool.REFRESH_SELECTED_WORKER)+1); "" != ips {
					//15秒钟没有返回结果则退出协程(截图比较耗时)
					ctxTimeout, _ := context.WithTimeout(this.ctx, time.Second*15)
					this.goroutine.Submit(func() {
						if messages, err := this.mill.Subscribe(ctxTimeout, tool.WORKER_ADMIN_CTX_RESULT); nil == err {
							//可能有多台主机,因此循环
							for {
								select {
								case msg := <-messages:
									msg.Ack()
									switch msg.Metadata.Get(tool.WORKER_ADMIN_CTX_SUB_TYPE) {
									case tool.WORKER_UPDATE:
										worker := model.Worker{}
										if nil == jsoniter.Unmarshal(msg.Payload, &worker) {
											if res, err := jsoniter.Marshal(Ping{Query: tool.WORKER_UPDATE, Response: worker}); nil == err {
												//请求需要时间(搞不好就关闭了)
												if !session.IsClosed() {
													session.Write(res)
												}
											}
										}
										break
									}
									break
								case <-ctxTimeout.Done():
									goto EXIT_UPDATE
								}
							}
						EXIT_UPDATE:
							//code here
						}
					})
					this.mill.Publish(tool.WORKER_ADMIN_CTX, ext.NewAdminSubTypeMessage([]byte(ips), tool.REFRESH_SELECTED_WORKER))
				}
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
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/v1?token=%s", token))
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
