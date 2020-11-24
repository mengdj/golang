package main

import (
	"broadcast"
	"context"
	"github.com/alecthomas/log4go"
	"github.com/panjf2000/ants/v2"
	"os"
	"os/signal"
	"server"
	"syscall"
	"web/model"
)

const (
	BROADCAST_PORT        = 10000
	BROADCAST_LISTEN_PORT = 10001
	WEB_LISTEN_PORT       = 10002
)

func main() {
	var (
		ants_pool *ants.Pool
		signChan  chan os.Signal
		err       error
	)
	logger := log4go.NewDefaultLogger(log4go.DEBUG)
	logger.LoadConfiguration("config/log4go.xml")
	ctx, cancel_ctx := context.WithCancel(context.Background())
	ants_pool, err = ants.NewPool(5, ants.WithPreAlloc(false))
	if nil != err {
		panic(err)
	}
	signChan = make(chan os.Signal)
	defer func() {
		close(signChan)
	}()
	signal.Notify(signChan, syscall.SIGINT, syscall.SIGTERM)
	bc := broadcast.NewBroadcast(&logger)
	app := server.NewApp(&logger, model.Port{BROADCAST_LISTEN_PORT,WEB_LISTEN_PORT})
	ants_pool.Submit(func() {
		if err := bc.Start(ctx, BROADCAST_LISTEN_PORT, BROADCAST_PORT); err != nil {
			logger.Error("broadcast: %s\n", err)
		}
	})
	ants_pool.Submit(func() {
		if err := app.Start(ctx); nil != err {
			logger.Error("server: %s\n", err)
		}
	})
	<-signChan
	logger.Info("关闭中...")
	cancel_ctx()
}
