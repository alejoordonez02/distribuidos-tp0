package common

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

const (
	BUF_SIZE     = 1024
	LEN_SIZE int = 2
)

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
	len := uint16(len(bytes_msg))

	binary.Write(bytes, binary.BigEndian, len)
	bytes.Write(bytes_msg)

	c.send(bytes.Bytes())

	return nil
}

func (c *Conn) Recv() (Response, error) {
	bytes := make([]byte, BUF_SIZE)
	err := c.readExact(bytes, LEN_SIZE)
	if err != nil {
		return Response{}, err
	}

	len_msg := int(binary.BigEndian.Uint16(bytes[:LEN_SIZE]))
	if LEN_SIZE+len_msg > BUF_SIZE {
		return Response{},
			fmt.Errorf(
				"message too big, size is %v and payload buf size is %v",
				len_msg, BUF_SIZE-LEN_SIZE)
	}

	err = c.readExact(bytes[LEN_SIZE:], len_msg)
	if err != nil {
		return Response{}, err
	}

	response, err := Deserialize(bytes[LEN_SIZE : LEN_SIZE+len_msg])
	if err != nil {
		return Response{}, err
	}

	return response, nil
}

func (c *Conn) Close() error {
	return c.skt.Close()
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

func (c *Conn) readExact(buf []byte, amount int) error {
	if amount == 0 {
		return nil
	}

	_, err := io.ReadFull(c.skt, buf[:amount])
	return err
}
