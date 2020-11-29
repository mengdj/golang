package broadcast

import (
	"context"
	"fmt"
	"github.com/alecthomas/log4go"
	"github.com/golang/protobuf/proto"
	cron "github.com/robfig/cron/v3"
	"net"
	"proc"
	"tool"
)

const (
	UDP               = "udp"
	TCP               = "tcp"
	UDP_BROADCAST_ASK = "BROADCAST"
)

type Broadcast struct {
	crontab *cron.Cron
	logger  *log4go.Logger
}

func NewBroadcast(logger *log4go.Logger) *Broadcast {
	return &Broadcast{crontab: cron.New(cron.WithSeconds()), logger: logger}
}

func (this *Broadcast) Start(ctx context.Context, listenPort, broadcastPort uint32) error {
	var err error = nil
	var conn net.Conn
	go func() {
		conn, err = net.Dial(UDP, fmt.Sprintf("255.255.255.255:%d", broadcastPort))
		if nil == err {
			var (
				pack       proc.Broadcast
				cmd        string = UDP_BROADCAST_ASK
				ip         string
				bodyLength uint32
				writeData  []byte
			)
			pack.Head = &proc.Head{}
			pack.Body = &proc.Body{}
			pack.Head.Cmd = &cmd
			//获取本机ip地址
			if ip, err = tool.LocalAddress(); nil == err {
				pack.Body.Server = &ip
				pack.Body.Port = &listenPort
				bodyLength = uint32(len([]byte(pack.Body.String())))
				pack.Head.Length = &bodyLength
				if writeData, err = proto.Marshal(&pack); nil == err {
					this.crontab.AddFunc("*/10 * * * * *", func() {
						_, err = conn.Write(writeData)
						if nil != err {
							this.logger.Warn(err)
						}
					})
					this.crontab.Start()
					this.logger.Info("发布广播:%d", broadcastPort)
				} else {
					this.logger.Critical(err)
				}
			}
		} else {
			this.logger.Critical(err)
		}
	}()
	<-ctx.Done()
	this.crontab.Stop()
	this.logger.Info("结束广播")
	return err
}
