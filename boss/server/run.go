package server

import (
	"bufio"
	"codec"
	"context"
	"fmt"
	"github.com/alecthomas/log4go"
	"github.com/gogf/gf/container/gpool"
	"github.com/gogf/gf/container/gtype"
	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
	"github.com/robfig/cron"
	"os"
	"proc"
	"strings"
	"time"
	"tool"
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
	addr      string
	multicore bool
	async     bool
	//协程池
	workerPool *goroutine.Pool
	//协议处理
	cmdPool *gpool.Pool
	//chan切片池
	chanBytePool *gpool.Pool
	logger       *log4go.Logger
	//在线连接
	connects ConnectContainer
	//定时任务
	cronTask *cron.Cron
}

func NewApp(logger *log4go.Logger) *App {
	return &App{logger: logger, multicore: true, async: true, connects: make(map[string]*connectItem), cronTask: cron.New()}
}

func (this *App) Start(ctx context.Context, listenPort uint32) error {
	//使用自定义解码器解码
	return gnet.Serve(this, fmt.Sprintf("%s://0.0.0.0:%d", TCP, listenPort), gnet.WithMulticore(this.multicore), gnet.WithTCPKeepAlive(time.Minute*5), gnet.WithCodec(codec.NewCodec()))
}

func (this *App) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	this.logger.Info("启动成功 %s (multi-cores: %t, loops: %d)",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
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
	//协程池
	this.workerPool = goroutine.Default()
	//每15秒输出一次已连接的客户端数并检测存活的连接
	this.cronTask.AddFunc("*/15 * * * *", func() {
		now := time.Now().Unix()
		active := 0
		for i, v := range this.connects {
			if (now - v.ping) > 15 {
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
		this.logger.Info("在线客户:%d", active)
	})
	this.cronTask.Start()
	return
}

func (this *App) Tick() (delay time.Duration, action gnet.Action) {
	return time.Second * 15, gnet.None
}

func (this *App) OnShutdown(svr gnet.Server) {
	this.cronTask.Stop()
	this.workerPool.Release()
	this.chanBytePool.Close()
	this.cmdPool.Close()
}

func (this *App) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	if action != gnet.Close && action != gnet.Shutdown {
		if len(frame) > 0 {
			remoteAddr := c.RemoteAddr().String()
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
	remoteAddr := c.RemoteAddr().String()
	this.connects[remoteAddr] = &connectItem{conn: c, close: make(chan None), ping: time.Now().Unix()}
	this.connects[remoteAddr].cap = capture{nil, nil, gtype.NewBool(true)}
	if i, e := this.chanBytePool.Get(); nil == e {
		if c, ok := i.(chan []byte); ok {
			this.connects[remoteAddr].cap.recv = c
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
	removeAddr := c.RemoteAddr().String()
	if s, ok := this.connects[removeAddr]; ok {
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