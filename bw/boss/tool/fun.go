package tool

import (
	"github.com/gogf/gf/util/grand"
	"net"
	"time"
)

//尝试调用某个可能panic的函数，如果出错就捕捉，避免退出
func Try(f func(), e func(i interface{})) {
	defer func() {
		if r := recover(); nil != r {
			e(r)
		}
	}()
	f()
}

func Now() int64{
	return time.Now().Unix()
}
//获取本机ip地址
func LocalAddress() (string, error) {
	var (
		addrs []net.Addr
		err   error
		ret   string = ""
	)
	addrs, err = net.InterfaceAddrs()
	if nil == err {
		for _, val := range addrs {
			if ipNet, ok := val.(*net.IPNet); ok {
				if nil != ipNet.IP.To4() {
					ret = ipNet.IP.String()
				}
			}
		}
	}
	return ret, err
}

func RandomStr(chars string, size int) string {
	b := make([]byte, size)
	for i := 0; i < size; i++ {
		b[i] = chars[grand.Intn(len(chars))]
	}
	return string(b)
}
