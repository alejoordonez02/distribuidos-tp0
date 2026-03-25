package messages

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/src/comms/protocol"
)

func (bets BetBatch) Serialize() []byte {
	var serial bytes.Buffer
	serial.WriteByte(protocol.TYPE_BET_BATCH)
	size := uint16(len(bets))
	binary.Write(&serial, binary.BigEndian, size)
	for _, b := range bets {
		b.serializeInto(&serial)
	}

	return serial.Bytes()
}

func (q Query) Serialize() []byte {
	serial := []byte{protocol.TYPE_QUERY}
	return serial
}

func Deserialize(bytes []byte) (Message, error) {
	switch bytes[0] {
	case protocol.TYPE_ACK:
		return deserializeAck(bytes[1:])
	case protocol.TYPE_NACK:
		return deserializeNack(bytes[1:])
	case protocol.TYPE_RESPONSE:
		return deserializeResponse(bytes[1:]), nil
	default:
		return Ack{}, errors.New("unknown response type")
	}
}

func deserializeAck(_ []byte) (Ack, error) {
	return Ack{true}, nil
}

func deserializeNack(_ []byte) (Ack, error) {
	return Ack{false}, nil
}

func deserializeResponse(serial []byte) Response {
	WinnerAmount := int(binary.BigEndian.Uint16(serial[:protocol.LEN_WINNER_AMOUNT]))
	response := Response{WinnerAmount}
	return response
}

func (b *Bet) serializeInto(buf *bytes.Buffer) []byte {
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
