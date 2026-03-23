package common

import (
	"bytes"
	"encoding/binary"
)

const TYPE_BET = 0x00

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

func (p *Person) serializeInto(buf *bytes.Buffer) {
	writeString(p.Name, buf)
	writeString(p.Surname, buf)
	binary.Write(buf, binary.BigEndian, p.birth)
}

func writeString(s string, buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, len(s))
	buf.WriteString(s)
}
