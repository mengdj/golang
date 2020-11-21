package server

import (
	"context"
	"fmt"
	"github.com/alecthomas/log4go"
	"github.com/gogf/gf/container/gpool"
	"github.com/golang/protobuf/proto"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
	"github.com/robfig/cron"
	"proc"
	"codec"
	"time"
)

const (
	TCP = "tcp"
)

type ConnectItem = struct {
	conn     gnet.Conn
	name     string
	closed   bool
	prevPing int64
}

type ConnectContainer = map[string]ConnectItem

type App struct {
	*gnet.EventServer
	addr       string
	multicore  bool
	async      bool
	codec      gnet.ICodec
	workerPool *goroutine.Pool
	cmdPool    *gpool.Pool
	logger     *log4go.Logger
	connects   ConnectContainer
	cronTask   *cron.Cron
}

func NewApp(logger *log4go.Logger) *App {
	return &App{logger: logger, multicore: false, async: true, connects: make(map[string]ConnectItem), cronTask: cron.New()}
}

func (this *App) Start(ctx context.Context, listenPort uint32) error {
	//使用自定义解码器解码
	return gnet.Serve(this, fmt.Sprintf("%s://0.0.0.0:%d", TCP, listenPort), gnet.WithMulticore(this.multicore), gnet.WithTCPKeepAlive(time.Minute*5), gnet.WithCodec(codec.NewCodec()))
}

func (this *App) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	this.logger.Info("服务器启动成功 %s (multi-cores: %t, loops: %d)",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	this.cmdPool = gpool.New(0, func() (interface{}, error) {
		ret := new(proc.Cmd)
		ret.Head = new(proc.Head)
		ret.Content = new(proc.Content)
		return ret, nil
	}, nil)
	this.workerPool = goroutine.Default()
	//每15秒输出一次已连接的客户端数并检测存活的连接
	this.cronTask.AddFunc("*/15 * * * *", func() {
		now := time.Now().Unix()
		active := 0
		for i, v := range this.connects {
			if (now - v.prevPing) > 30 {
				if err := v.conn.Close(); nil != err {
					this.logger.Warn(err)
				}
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

func (this *App) OnShutdown(svr gnet.Server) {
	this.workerPool.Release()
}

func (this *App) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	if action != gnet.Close && action != gnet.Shutdown {
		if len(frame) > 0 {
			remoteAddr := c.RemoteAddr().String()
			//read
			if cmd,err := this.cmdPool.Get();nil==err{
				if nil != cmd {
					msg := cmd.(*proc.Cmd)
					if err := proto.Unmarshal(frame, msg); nil == err {
						switch *msg.Content.Type {
						case proc.ContentType_PING:
							//处理PING请求(用ping请求来辅助检测客户端是否还存活)
							if s, ok := this.connects[remoteAddr]; ok {
								s.prevPing = time.Now().Unix()
								s.name = *msg.Content.GetPing().Name
								this.connects[remoteAddr] = s
							}
							break
						case proc.ContentType_CAPTURE:
							this.logger.Info(*(msg.Content.Param.(*proc.Content_Capture).Capture.Id),*(msg.Content.Param.(*proc.Content_Capture).Capture.Seq))
							break
						}
					} else {
						this.logger.Info("数据包异常:%s",err.Error())
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
	this.connects[remoteAddr] = ConnectItem{conn: c, closed: false, prevPing: time.Now().Unix()}
	return
}

// OnClosed fires when a connection has been closed.
// The parameter:err is the last known connection error.
func (this *App) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	removeAddr := c.RemoteAddr().String()
	if s, ok := this.connects[removeAddr]; ok {
		if nil != s.conn {
			delete(this.connects, removeAddr)
		}
	}
	return
}

func (this *App) Stop() error {
	this.cronTask.Stop()
	return nil
}
