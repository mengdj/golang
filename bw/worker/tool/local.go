package tool

import (
	"github.com/shirou/gopsutil/host"
	"golang.org/x/sync/singleflight"
)

type Local struct {
	singleSetCache singleflight.Group
}

func NewLocal() *Local {
	return &Local{}
}

func (this *Local) InfoStat() (*host.InfoStat, error) {
	ret, err, _ := this.singleSetCache.Do("Local_Info", func() (interface{}, error) {
		return host.Info()
	})
	if nil != err {
		return nil, err
	}
	return ret.(*host.InfoStat), err
}
