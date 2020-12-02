package tool

import (
	"github.com/gogf/gf/util/grand"
	"net"
	"time"
)

//尝试调用某个可能panic的函数，如果出错就捕捉，避免退出(此函数可能导致性能损失)
func Try(f func(), e func(i interface{})) {
	defer func() {
		if r := recover(); nil != r {
			e(r)
		}
	}()
	f()
}

func If(ok bool,f func()){
	if ok{
		f()
	}
}

func Now() int64 {
	return time.Now().Unix()
}

//获取本机ip地址
func LocalAddress() (string, error) {
	var (
		addrs      []net.Addr
		interfaces []net.Interface
		err        error
		ret        string = ""
	)
	interfaces, err = net.Interfaces()
	if nil != err {
		return ret, err
	}
GET_ADDR:
	for _, f := range interfaces {
		//可能多块网卡。。。
		if 0 != (f.Flags & net.FlagUp) {
			addrs, err = f.Addrs()
			if nil != err {
				return ret, err
			}
			for _, addr := range addrs {
				if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
					if nil != ipNet.IP.To4() {
						ret = ipNet.IP.String()
						break GET_ADDR
					}
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
