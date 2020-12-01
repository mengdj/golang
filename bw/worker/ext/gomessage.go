package ext

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gogf/gf/container/gpool"
)

var globalMessagePool *gpool.Pool=nil

type ExtMessage struct{
	*message.Message
}

func (this *ExtMessage) Ack() bool{
	return this.Message.Ack()
}

func init() {
	//初始化对象池
	globalMessagePool=gpool.New(0, func() (interface{}, error) {
		return message.NewMessage(watermill.NewUUID(),[]byte{}),nil
	}, func(i interface{}) {
	})
}


