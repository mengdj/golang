package server

import (
	"bufio"
	"codec"
	"context"
	"encoding/json"
	"ext"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/alecthomas/log4go"
	"github.com/gogf/gf/container/gpool"
	"github.com/gogf/gf/container/gtype"
	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
	"os"
	"proc"
	"strings"
	"time"
	"tool"
	www "web"
	"web/model"
)

const (
	TCP = "tcp"
)

type None = struct{}
type capture struct {
	recv   chan []byte
	file   *os.File
	status *gtype.Bool
}
type connectItem = struct {
	//连接
	conn gnet.Conn
	//客户端机器名称
	name string
	//捕捉数据接受channel
	cap capture
	//客户端断开channel
	close chan None
	//上一次ping的事假
	ping int64
}

type ConnectContainer = map[string]*connectItem

type App struct {
	*gnet.EventServer
	port model.Port
	//协程池
	workerPool *goroutine.Pool
	//协议处理
	cmdPool *gpool.Pool
	//chan切片池
	chanBytePool *gpool.Pool
	logger       *log4go.Logger
	//在线连接
	connects ConnectContainer
	//web服务器
	web  *www.Web
	mill *ext.ExtGoChanel
}

func NewApp(logger *log4go.Logger, p model.Port) *App {
	return &App{logger: logger, connects: make(map[string]*connectItem), port: p}
}

/** 使用自定义解码器解码 */
func (this *App) Start(ctx context.Context) error {
	return gnet.Serve(this, fmt.Sprintf("%s://0.0.0.0:%d", TCP, this.port.Socket), gnet.WithLogger(Logger{this.logger}), gnet.WithTicker(true), gnet.WithReusePort(true), gnet.WithMulticore(true), gnet.WithTCPKeepAlive(time.Minute*1), gnet.WithCodec(codec.NewCodec()))
}

func (this *App) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	//协程池
	this.workerPool = goroutine.Default()
	this.logger.Info("启动成功 %s (multi-cores: %t, loops: %d)", srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	this.mill = ext.NewExtGoChannel(gochannel.Config{Persistent: true, BlockPublishUntilSubscriberAck: false}, watermill.NewStdLogger(false, false))
	this.web = www.NewWeb(this.logger, this.mill, this.workerPool)
	//协议池
	this.cmdPool = gpool.New(0, func() (interface{}, error) {
		ret := new(proc.Cmd)
		ret.Head = new(proc.Head)
		ret.Content = new(proc.Content)
		return ret, nil
	}, nil)
	//切片池
	this.chanBytePool = gpool.New(0, func() (interface{}, error) {
		return make(chan []byte), nil
	}, func(i interface{}) {
		if c, ok := i.(chan []byte); ok {
			close(c)
		}
	})
	//启动web服务器
	this.workerPool.Submit(func() {
		if err := this.web.Run(fmt.Sprintf(":%d", this.port.Web)); nil != err {
			this.logger.Critical(err)
		}
	})
	this.workerPool.Submit(func() {
		//订阅查询所有客户端消息，给其他查询
		if messages, err := this.mill.Subscribe(context.Background(), tool.QUERY_WORKERS); nil == err {
			for {
				select {
				case <-messages:
					//查询所在的客户信息(传输到消息队列即可)
					workers := []model.Worker{}
					for s, c := range this.connects {
						workers = append(workers, model.Worker{Addr: s, Name: c.name, Ping: c.ping})
					}
					if res, err := json.Marshal(workers); nil == err {
						this.mill.Publish(tool.QUERY_WORKERS_RESULT, message.NewMessage(watermill.NewUUID(), res))
					}
				}
			}
		}
	})
	return
}

func (this *App) Tick() (delay time.Duration, action gnet.Action) {
	now := time.Now().Unix()
	active := 0
	for i, v := range this.connects {
		if (now - v.ping) > 10 {
			if err := v.conn.Close(); nil != err {
				this.logger.Warn(err)
			}
			v.close <- None{}
			this.chanBytePool.Put(v.cap.recv)
			delete(this.connects, i)
			continue
		}
		active++
	}
	//触发消息服务
	if active == 0 {
		return time.Second * 10, gnet.None
	}
	return time.Second * 5, gnet.None
}

func (this *App) OnShutdown(svr gnet.Server) {
	this.workerPool.Release()
	this.chanBytePool.Close()
	this.cmdPool.Close()
}

func (this *App) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	if action != gnet.Close && action != gnet.Shutdown {
		if len(frame) > 0 {
			remoteAddr := strings.Split(c.RemoteAddr().String(), ":")[0]
			//read
			if cmd, err := this.cmdPool.Get(); nil == err {
				if nil != cmd {
					msg := cmd.(*proc.Cmd)
					if err := proto.Unmarshal(frame, msg); nil == err {
						switch *msg.Content.Type {
						case proc.ContentType_PING:
							//处理PING请求(用ping请求来辅助检测客户端是否还存活)
							if s, ok := this.connects[remoteAddr]; ok {
								s.ping = time.Now().Unix()
								s.name = *msg.Content.GetPing().Name
								this.connects[remoteAddr] = s
							}
							break
						case proc.ContentType_CAPTURE:
							//发送数据给接收者接收文件数据(状态异常就丢弃数据了，防止都塞)
							if this.connects[remoteAddr].cap.status.Val() {
								if *(msg.Content.GetCapture().Compress) {
									if tmp, derr := snappy.Decode(nil, msg.Content.GetCapture().GetData()); nil == derr {
										this.connects[remoteAddr].cap.recv <- tmp
									} else {
										this.logger.Error(derr)
										//恢复一个异常解码数据
									}
								} else {
									this.connects[remoteAddr].cap.recv <- msg.Content.GetCapture().GetData()
								}
							}
							break
						}
					} else {
						this.logger.Info("数据包异常:%s", err.Error())
					}
					this.cmdPool.Put(cmd)
				}
			}
		}
	}
	return
}

func (this *App) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	remoteAddr := strings.Split(c.RemoteAddr().String(), ":")[0]
	nowUnix := time.Now().Unix()
	this.connects[remoteAddr] = &connectItem{conn: c, close: make(chan None), ping: nowUnix}
	this.connects[remoteAddr].cap = capture{nil, nil, gtype.NewBool(true)}
	if i, e := this.chanBytePool.Get(); nil == e {
		if c, ok := i.(chan []byte); ok {
			this.connects[remoteAddr].cap.recv = c
		}
	}
	//发布客户端上线消息(保证至少一个管理员上线的情况下才发布消息)
	if !this.mill.IsPause(tool.WORKER_ONLINE) {
		if d, e := json.Marshal(model.Worker{Addr: remoteAddr, Name: "", Ping: nowUnix}); nil == e {
			this.mill.Publish(tool.WORKER_ONLINE_RESULT, message.NewMessage(watermill.NewUUID(), d))
		}
	}
	//接收截图
	this.workerPool.Submit(func() {
		//切割ip数据
		if dir, e := os.Getwd(); nil == e {
			dir = fmt.Sprintf("%s/data/%s", dir, strings.Split(remoteAddr, ":")[0])
			if e = os.MkdirAll(dir, os.ModePerm); e != nil {
				this.logger.Error(e)
				this.connects[remoteAddr].cap.status.Set(false)
			} else {
				//目标文件路径
				save := fmt.Sprintf("%s/capture.jpg", dir)
				//缓冲区写入文件
				var writer *bufio.Writer = nil
				for {
					if c, ok := this.connects[remoteAddr]; ok {
						select {
						case <-c.close:
							goto CLOSE_CON
						case data := <-c.cap.recv:
							//code here
							if c.cap.status.Val() {
								if 0 == len(data) {
									if nil == c.cap.file {
										//start
										if c.cap.file, e = os.OpenFile(save, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm); nil == e {
											if nil == writer {
												writer = bufio.NewWriter(c.cap.file)
											}
											writer.Reset(c.cap.file)
										} else {
											this.logger.Error(e)
											c.cap.status.Set(false)
										}
									} else {
										//stop
										writer.Flush()
										c.cap.file.Close()
										c.cap.file = nil
									}
								} else {
									if nil != c.cap.file && nil != writer {
										writer.Write(data)
									}
								}
							}
						}
					}
				}
			CLOSE_CON:
				//
			}
		}
	})
	return
}

func (this *App) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	removeAddr := strings.Split(c.RemoteAddr().String(), ":")[0]
	if s, ok := this.connects[removeAddr]; ok {
		//发布客户端离线消息
		if !this.mill.IsPause(tool.WORKER_OFFLINE) {
			if d, e := json.Marshal(model.Worker{Addr: removeAddr, Name: s.name, Ping: s.ping}); nil == e {
				this.mill.Publish(tool.WORKER_OFFLINE_RESULT, message.NewMessage(watermill.NewUUID(), d))
			}
		}
		if nil != s.conn {
			s.close <- None{}
			tool.Try(func() {
				close(s.close)
			}, func(i interface{}) {
				this.logger.Info(i)
			})
			this.chanBytePool.Put(s.cap.recv)
			delete(this.connects, removeAddr)
		}
	}
	return
}
