package protocol84

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

type LoginPacket struct {
	Protocol      uint32
	Username      string
	UUID          uuid.UUID
	ClientId      uint64
	PublicKey     string
	ServerAddress string
}

func (packet *LoginPacket) decodeTokenToObject(token string) (map[string]interface{}, error) {
	var data map[string]interface{}

	jsonPayload, err := base64.StdEncoding.DecodeString(token)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonPayload, &data)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (packet *LoginPacket) readChainData(chainData map[string]interface{}) error {
	var chain []string

	if chainData["chain"] == nil {
		return errors.New("the chain data doesn't have a 'chain' key")
	}

	chain = chainData["chain"].([]string)

	for _, token := range chain {
		var _tokenData map[string]interface{}

		_tokenData, err := packet.decodeTokenToObject(token)

		if err != nil {
			continue
		}
	}

	return nil
}

func (packet *LoginPacket) Read(buffer *bytes.Buffer) error {
	var err error
	var zlibBufferSize uint32

	err = binary.Read(buffer, binary.BigEndian, &packet.Protocol)

	if err != nil {
		return err
	}

	err = binary.Read(buffer, binary.BigEndian, &zlibBufferSize)

	if err != nil {
		return err
	}

	zlibBuffer := make([]byte, zlibBufferSize)
	reader, err := zlib.NewReader(buffer)

	if err != nil {
		return err
	}

	_, err = reader.Read(zlibBuffer)

	if err != nil {
		return err
	}

	var chainData map[string]interface{}

	err = json.Unmarshal(zlibBuffer, &chainData)

	if err != nil {
		return err
	}

	return packet.readChainData(chainData)
}
