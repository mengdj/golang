package main

import (
	"admin/model"
	"broadcast"
	"context"
	"ext"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/alecthomas/log4go"
	"github.com/panjf2000/ants/v2"
	"os"
	"os/signal"
	"server"
	"sync"
	"syscall"
)

const (
	//组播端口
	BROADCAST_PORT        = 10000
	//tcp服务端口
	BROADCAST_LISTEN_PORT = 10001
	//web管理端口
	WEB_LISTEN_PORT       = 10002
)

func main() {
	var (
		goroutine *ants.Pool
		err       error
		multiTask sync.WaitGroup
	)
	logger := log4go.NewDefaultLogger(log4go.DEBUG)
	logger.LoadConfiguration("config/log4go.xml")
	pubSub := ext.NewExtGoChannel(gochannel.Config{Persistent: false, BlockPublishUntilSubscriberAck: false}, watermill.NewStdLogger(false, false))
	ctx, cancelCtx := context.WithCancel(context.Background())
	goroutine, err = ants.NewPool(2, ants.WithPreAlloc(false))
	//遵循 一个协程不知道停止它，就不要创建它的原则 避免内存泄漏
	signTermChan := make(chan os.Signal)
	defer func() {
		goroutine.Release()
		close(signTermChan)
		pubSub.Close()
		logger.Close()
	}()
	if nil != err {
		panic(err)
	}
	signal.Notify(signTermChan, syscall.SIGINT, syscall.SIGTERM)
	serverBroadcast := broadcast.NewBroadcast(&logger)
	serverApplication := server.NewApp(ctx, pubSub, &logger, model.Port{BROADCAST_LISTEN_PORT, WEB_LISTEN_PORT})
	goroutine.Submit(func() {
		multiTask.Add(1)
		if err := serverBroadcast.Start(ctx, BROADCAST_LISTEN_PORT, BROADCAST_PORT); err != nil {
			logger.Error("broadcast: %s\n", err)
		}
		multiTask.Done()
	})
	goroutine.Submit(func() {
		multiTask.Add(1)
		if err := serverApplication.Start(ctx); nil != err {
			logger.Error("server: %s\n", err)
		}
		multiTask.Done()
	})
	<-signTermChan
	cancelCtx()
	multiTask.Wait()
	logger.Info("服务器已关闭，感谢您的使用")
}
