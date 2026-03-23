package common

import (
	"bytes"
	"encoding/binary"
	"net"
)

const LEN_SIZE = 2

type Conn struct {
	skt net.Conn
}

func NewConn(addr string) (Conn, error) {
	skt, err := net.Dial("tcp", addr)
	if err != nil {
		return Conn{}, err
	}

	conn := Conn{skt}
	return conn, nil
}

func (c *Conn) Send(msg Serializable) error {
	bytes := new(bytes.Buffer)
	bytes_msg := msg.Serialize()
	len := len(bytes_msg)

	binary.Write(bytes, binary.BigEndian, len)
	bytes.Write(msg.Serialize())

	c.send(bytes.Bytes())

	return nil
}

func (c *Conn) Close() error {
	err := c.skt.Close()
	if err != nil {
		return err
	}

	return nil
}

func (c *Conn) send(bytes []byte) error {
	sent := 0
	for sent < len(bytes) {
		_sent, err := c.skt.Write(bytes[sent:])
		if err != nil {
			return err
		}

		sent += _sent
	}

	return nil
}
