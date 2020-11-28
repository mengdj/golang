package codec

import (
	"encoding/binary"
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
