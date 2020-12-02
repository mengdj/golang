package broadcast

import (
	"context"
	"ext"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/alecthomas/log4go"
	"net"
)

const (
	UDP               = "udp"
	TCP               = "tcp"
	UDP_BROADCAST_ASK = "BROADCAST"
)

type Broadcast struct {
	running bool
	logger  *log4go.Logger
}

func NewBroadcast(logger *log4go.Logger) *Broadcast {
	return &Broadcast{running: false, logger: logger}
}

func (this *Broadcast) Start(ctx context.Context, listenPort uint32, pub *ext.ExtGoChanel) error {
	var (
		conn *net.UDPConn = nil
		err  error        = nil
		addr *net.UDPAddr
	)
	defer func() {
		if nil != conn {
			conn.Close()
		}
	}()
	if !this.running {
		addr, err = net.ResolveUDPAddr(UDP, fmt.Sprintf("0.0.0.0:%d", listenPort))
		if nil != err {
			this.logger.Error(err)
		}
		conn, err = net.ListenUDP(UDP, addr)
		if nil != err {
			this.logger.Error(err)
		}
		this.logger.Info("接收广播:%d", listenPort)
		this.running = true
		go func() {
			for {
				var buffer [512]byte
				retCount, _, retErr := conn.ReadFromUDP(buffer[0:])
				if nil == retErr &&retCount > 0 {
					if !pub.IsPause("START_CONNECT_SERVER") {
						msg := message.NewMessage(watermill.NewUUID(), buffer[:retCount])
						if err := pub.Publish("START_CONNECT_SERVER", msg); err != nil {
							this.logger.Warn(err)
						}
					}
				}else{
					this.logger.Critical(retErr)
				}
			}
		}()
		//wait
		<-ctx.Done()
		this.running = false
	}
	return err
}
