package codec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/gogf/gf/net/gtcp"
	"net"
)

type Codec struct {
	*gtcp.Conn
}

const (
	PACKAGE_TAG string = "CMD"
)

func NewCodec(con net.Conn) *Codec {
	return &Codec{Conn: gtcp.NewConnByNetConn(con)}
}

//写入应用协议protobuff
func (this *Codec) Write(data []byte) error {
	if size := len(data); size > 0 {
		buffer := bytes.NewBuffer([]byte{})
		buffer.WriteString(PACKAGE_TAG)
		//数据包的大小，不包含包头
		binary.Write(buffer, binary.BigEndian, uint32(size))
		buffer.Write(data)
		return this.Send(buffer.Bytes())
	}
	return errors.New("数据包不能为空")
}

func (this *Codec) Read() (ret []byte, err error) {
	var buffer []byte
	var length int
	// cmd+4=7
	buffer, err = this.Recv(7)
	if err != nil || string(buffer[:3]) != PACKAGE_TAG {
		//检验接受包的数据
		return nil, err
	}
	length = int(binary.BigEndian.Uint32([]byte{buffer[3], buffer[4], buffer[5], buffer[6]}))
	if length < 0 {
		return nil, fmt.Errorf(`invalid package size %d`, length)
	}else if length == 0 {
		return nil, nil
	}
	return this.Recv(length)
}
