package main

import (
	"broadcast"
	"client"
	"context"
	"ext"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/alecthomas/log4go"
	"github.com/golang/protobuf/proto"
	"os"
	"os/signal"
	"proc"
	"sync"
	"sync/atomic"
	"syscall"
)

const (
	UDP               = "udp"
	TCP               = "tcp"
	UDP_BROADCAST_ASK = "BROADCAST_ASK"
	BROADCAST_PORT    = 10000
)

func main() {
	var (
		sign         chan os.Signal
		broadcastLis *broadcast.Broadcast
		cli          *client.Client
		useCount     uint32
	)
	broadcastPool := sync.Pool{New: func() interface{} {
		ret := new(proc.Broadcast)
		ret.Head = new(proc.Head)
		ret.Body = new(proc.Body)
		return ret
	}}
	logger := log4go.NewDefaultLogger(log4go.DEBUG)
	logger.LoadConfiguration("config/log4go.xml")
	pubSub := ext.NewExtGoChannel(gochannel.Config{Persistent: true, BlockPublishUntilSubscriberAck: true}, watermill.NewStdLogger(false, false))
	messages, err := pubSub.Subscribe(context.Background(), "START_CONNECT_SERVER")
	if nil != err {
		panic(err)
	}
	broadcastLis = broadcast.NewBroadcast(&logger)
	cli = client.NewClient(&logger, pubSub)
	sign = make(chan os.Signal)
	defer func() {
		close(sign)
	}()
	ctx, cancel_ctx := context.WithCancel(context.Background())
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	//接收广播通知并告知创建客户端长连接(watermill)
	go func(messages <-chan *message.Message) {
		for msg := range messages {
			msg.Ack()
			if 0 == atomic.LoadUint32(&useCount) {
				atomic.AddUint32(&useCount, 1)
				if ins := broadcastPool.Get(); nil != ins {
					cast := ins.(*proc.Broadcast)
					cast.Reset()
					castErr := proto.Unmarshal(msg.Payload, cast)
					if nil == castErr {
						//开启客户端TCP连接
						pubSub.Pause("START_CONNECT_SERVER")
						if e := cli.Start(ctx, cast); nil != e {
							logger.Critical(e)
						}
						pubSub.Continue("START_CONNECT_SERVER")
					} else {
						logger.Debug(castErr)
					}
					broadcastPool.Put(ins)
				} else {
					logger.Warn("对象获取失败")
				}
				atomic.StoreUint32(&useCount, 0)
			}
		}
	}(messages)
	go func() {
		if e := broadcastLis.Start(ctx, BROADCAST_PORT, pubSub); nil != e {
			logger.Warn(e)
		}
	}()
	//获取中断事件
	<-sign
	cancel_ctx()
	pubSub.Close()
}
