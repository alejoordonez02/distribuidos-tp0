package common

import (
	"bytes"
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
	var serial bytes.Buffer
	serial.WriteByte(TYPE_BET)
	serializeStringInto(b.Agency, &serial)
	serializeStringInto(b.FirstName, &serial)
	serializeStringInto(b.LastName, &serial)
	serializeStringInto(b.Document, &serial)
	serializeStringInto(b.BirthDate, &serial)
	serializeStringInto(b.Number, &serial)

	return serial.Bytes()
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

// Serialize string into bytes and write them into the buffer,
// prepending one byte with the length of the string bytes
func serializeStringInto(s string, buf *bytes.Buffer) {
	serial_len := byte(uint8(len(s)))
	buf.WriteByte(serial_len)
	buf.WriteString(s)
}
