package ext

import (
	"errors"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/gogf/gf/container/gset"
)

type ExtGoChanel struct {
	*gochannel.GoChannel
	pauseTopics *gset.StrSet
}

func NewExtGoChannel(config gochannel.Config, logger watermill.LoggerAdapter) *ExtGoChanel {
	t := &ExtGoChanel{gochannel.NewGoChannel(config, logger), gset.NewStrSet()}
	return t
}

func (this *ExtGoChanel) Publish(topic string, messages ...*message.Message) error {
	if !this.IsPause(topic){
		return this.GoChannel.Publish(topic,messages...)
	}
	return errors.New(fmt.Sprintf("%s has pause.",topic))
}

func (this *ExtGoChanel) Pause(topics ...string) {
	for _,t:=range topics{
		this.pauseTopics.AddIfNotExist(t)
	}
}

func (this *ExtGoChanel) Continue(topics ...string) {
	for _,t:=range topics{
		this.pauseTopics.Remove(t)
	}
}

func (this *ExtGoChanel) IsPause(topic string) bool {
	return this.pauseTopics.Contains(topic)
}


