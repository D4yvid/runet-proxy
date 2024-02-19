package message

import (
	"bytes"
	"encoding/binary"
)

type ConnectedPing struct {
	PingID int64
}

func (pk *ConnectedPing) Write(buf *bytes.Buffer) {
	_ = binary.Write(buf, binary.BigEndian, IDConnectedPing)
	_ = binary.Write(buf, binary.BigEndian, pk.PingID)
}

func (pk *ConnectedPing) Read(buf *bytes.Buffer) error {
	return binary.Read(buf, binary.BigEndian, &pk.PingID)
}
