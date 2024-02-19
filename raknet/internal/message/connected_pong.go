package message

import (
	"bytes"
	"encoding/binary"
)

type ConnectedPong struct {
	PingId int64
}

func (pk *ConnectedPong) Write(buf *bytes.Buffer) {
	_ = binary.Write(buf, binary.BigEndian, IDConnectedPong)
	_ = binary.Write(buf, binary.BigEndian, pk.PingId)
}

func (pk *ConnectedPong) Read(buf *bytes.Buffer) error {
	return binary.Read(buf, binary.BigEndian, &pk.PingId)
}
