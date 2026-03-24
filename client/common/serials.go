package common

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const (
	TYPE_BET_BATCH     = 0x00
	TYPE_RESPONSE_ACK  = 0x01
	TYPE_RESPONSE_NACK = 0x02
)

type Serializable interface {
	Serialize() []byte
}

func (bets BetBatch) Serialize() []byte {
	var serial bytes.Buffer
	serial.WriteByte(TYPE_BET_BATCH)
	size := uint16(len(bets))
	binary.Write(&serial, binary.BigEndian, size)
	for _, b := range bets {
		b.serializeInto(&serial)
	}

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

func (b *Bet) serializeInto(buf *bytes.Buffer) []byte {
	// buf.WriteByte(TYPE_BET)
	serializeStringInto(b.Agency, buf)
	serializeStringInto(b.FirstName, buf)
	serializeStringInto(b.LastName, buf)
	serializeStringInto(b.Document, buf)
	serializeStringInto(b.BirthDate, buf)
	serializeStringInto(b.Number, buf)

	return buf.Bytes()
}

// Serialize string into bytes and write them into the buffer,
// prepending one byte with the length of the string bytes
func serializeStringInto(s string, buf *bytes.Buffer) {
	serial_len := byte(uint8(len(s)))
	buf.WriteByte(serial_len)
	buf.WriteString(s)
}
