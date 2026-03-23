package common

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const (
	TYPE_BET           = 0x00
	TYPE_RESPONSE_ACK  = 0x01
	TYPE_RESPONSE_NACK = 0x02
)

type Serializable interface {
	Serialize() []byte
}

func (b *Bet) Serialize() []byte {
	bytes := new(bytes.Buffer)
	binary.Write(bytes, binary.BigEndian, TYPE_BET)
	binary.Write(bytes, binary.BigEndian, b.Number)
	b.Person.serializeInto(bytes)

	return bytes.Bytes()
}

func Deserialize(bytes []byte) (Response, error) {
	switch bytes[0] {
	case TYPE_RESPONSE_ACK:
		return Response{true}, nil
	case TYPE_RESPONSE_NACK:
		return Response{false}, nil
	default:
		return Response{}, errors.New("unknown response type")
	}
}

func (p *Person) serializeInto(buf *bytes.Buffer) {
	writeString(p.Name, buf)
	writeString(p.Surname, buf)
	binary.Write(buf, binary.BigEndian, p.birth)
}

func writeString(s string, buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, len(s))
	buf.WriteString(s)
}
