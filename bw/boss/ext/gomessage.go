package ext

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"tool"
)

type Meta = string

func NewMessage(payload message.Payload, meta ...Meta) *message.Message {
	msg := message.NewMessage(watermill.NewUUID(), payload)
	//组合元数据(请成对传递)
	if size := len(meta); size > 0 && 0 == size%2 {
		key:=""
		for i,v:=range meta{
			if 0==i%2{
				key=v
			}else{
				msg.Metadata.Set(key,v)
			}
		}
	}
	return msg
}

func NewAdminSubTypeMessage(payload message.Payload,val Meta) *message.Message{
	return NewMessage(payload,tool.WORKER_ADMIN_CTX_SUB_TYPE,val)
}