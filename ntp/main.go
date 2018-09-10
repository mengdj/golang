package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	var (
		ntp    *Ntp
		buffer []byte
		err    error
		ret    int
	)
	//链接阿里云NTP服务器,NTP有很多免费服务器可以使用time.windows.com
	conn, err := net.Dial("udp", "ntp1.aliyun.com:123")
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
		conn.Close()
	}()
	ntp = NewNtp()
	conn.Write(ntp.GetBytes())
	buffer = make([]byte, 2048)
	ret, err = conn.Read(buffer)
	if err == nil {
		if ret > 0 {
			ntp.Parse(buffer, true)
			fmt.Println(fmt.Sprintf(
				"LI:%d\r\n版本:%d\r\n模式:%d\r\n精度:%d\r\n轮询:%d\r\n系统精度:%d\r\n延时:%ds\r\n最大误差:%d\r\n时钟表示:%d\r\n时间戳:%d %d %d %d\r\n",
				ntp.Li,
				ntp.Vn,
				ntp.Mode,
				ntp.Stratum,
				ntp.Poll,
				ntp.Precision,
				ntp.RootDelay,
				ntp.RootDispersion,
				ntp.ReferenceIdentifier,
				ntp.ReferenceTimestamp,
				ntp.OriginateTimestamp,
				ntp.ReceiveTimestamp,
				ntp.TransmitTimestamp,
			))
		}
	}
}
