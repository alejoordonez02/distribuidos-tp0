package common

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const (
	TYPE_BET_BATCH = 0x00
	TYPE_ACK       = 0x01
	TYPE_NACK      = 0x02
	TYPE_QUERY     = 0x03
	TYPE_RESPONSE  = 0x04

	LEN_WINNER_AMOUNT = 2
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

func (q Query) Serialize() []byte {
	serial := []byte{TYPE_QUERY}
	return serial
}

func Deserialize(bytes []byte) (Message, error) {
	switch bytes[0] {
	case TYPE_ACK:
		return Ack{true}, nil
	case TYPE_NACK:
		return Ack{false}, nil
	case TYPE_RESPONSE:
		return deserializeResponse(bytes[1:]), nil
	default:
		return Ack{}, errors.New("unknown response type")
	}
}

func (b *Bet) serializeInto(buf *bytes.Buffer) []byte {
	// buf.WriteByte(TYPE_BET)
	// serializeStringInto(b.Agency, buf)
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

func deserializeResponse(serial []byte) Response {
	WinnerAmount := int(binary.BigEndian.Uint16(serial[:LEN_WINNER_AMOUNT]))
	response := Response{WinnerAmount}
	return response
}
