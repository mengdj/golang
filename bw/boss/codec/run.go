package codec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/panjf2000/gnet"
)

const (
	PACKAGE_TAG string = "CMD"
)

var packageTagLength = len(PACKAGE_TAG)

type Codec struct {
	*gnet.LengthFieldBasedFrameCodec
}

func NewCodec() *Codec {
	return &Codec{gnet.NewLengthFieldBasedFrameCodec(gnet.EncoderConfig{binary.BigEndian, 4, -3, false}, gnet.DecoderConfig{binary.BigEndian, packageTagLength, 4, 0, 7})}
}

func (cc *Codec) Decode(c gnet.Conn) ([]byte, error){
	return cc.LengthFieldBasedFrameCodec.Decode(c)
}

//写入包头+长度+具体的内容
func (cc *Codec) Encode(c gnet.Conn, data []byte) ([]byte, error) {
	length := len(data)
	if length < 0 {
		return nil, errors.New("buf can't null")
	}
	buffer := bytes.NewBuffer([]byte{})
	buffer.Write([]byte(PACKAGE_TAG))
	binary.Write(buffer, binary.BigEndian, uint32(length))
	buffer.Write(data)
	return buffer.Bytes(), nil
}
