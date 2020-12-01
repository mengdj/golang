package main

import (
	"admin/model"
	"broadcast"
	"context"
	"github.com/alecthomas/log4go"
	"github.com/panjf2000/ants/v2"
	"os"
	"os/signal"
	"server"
	"sync"
	"syscall"
)

const (
	BROADCAST_PORT        = 10000
	BROADCAST_LISTEN_PORT = 10001
	WEB_LISTEN_PORT       = 10002
)

func main() {
	var (
		antsPool *ants.Pool
		signChan chan os.Signal
		err      error
		waitTask sync.WaitGroup
	)
	logger := log4go.NewDefaultLogger(log4go.DEBUG)
	logger.LoadConfiguration("config/log4go.xml")
	ctx, cancelCtx := context.WithCancel(context.Background())
	antsPool, err = ants.NewPool(5, ants.WithPreAlloc(false))
	if nil != err {
		panic(err)
	}
	//遵循 一个协程不知道停止它，就不要创建它的原则 避免内存泄漏
	signChan = make(chan os.Signal)
	defer func() {
		close(signChan)
		logger.Close()
	}()
	signal.Notify(signChan, syscall.SIGINT, syscall.SIGTERM)
	borad := broadcast.NewBroadcast(&logger)
	applic := server.NewApp(ctx, &logger, model.Port{BROADCAST_LISTEN_PORT, WEB_LISTEN_PORT})
	antsPool.Submit(func() {
		waitTask.Add(1)
		if err := borad.Start(ctx, BROADCAST_LISTEN_PORT, BROADCAST_PORT); err != nil {
			logger.Error("broadcast: %s\n", err)
		}
		waitTask.Done()
	})
	antsPool.Submit(func() {
		waitTask.Add(1)
		if err := applic.Start(ctx); nil != err {
			logger.Error("server: %s\n", err)
		}
		waitTask.Done()
	})
	<-signChan
	cancelCtx()
	//等待线下资源结束
	waitTask.Wait()
	antsPool.Release()
	logger.Info("服务器已关闭，感谢您的使用")
}
