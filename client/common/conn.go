package common

import (
	"bytes"
	"encoding/binary"
	"errors"
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
	len := len(bytes_msg)

	binary.Write(bytes, binary.BigEndian, len)
	bytes.Write(msg.Serialize())

	c.send(bytes.Bytes())

	return nil
}

func (c *Conn) Recv() (Response, error) {
	var bytes []byte
	read, err := c.readAtLeast(bytes, LEN_SIZE)
	if err != nil {
		return Response{}, err
	}

	len_msg := int(binary.BigEndian.Uint32(bytes[:LEN_SIZE]))
	missing := len_msg + LEN_SIZE - read
	read, err = c.readAtLeast(bytes[read:], missing)
	if err != nil {
		return Response{}, err
	}

	response, err := Deserialize(bytes)
	if err != nil {
		return Response{}, err
	}

	return response, nil
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

func (c *Conn) readAtLeast(buf []byte, atLeast int) (int, error) {
	read := 0
	for read < atLeast {
		_read, err := c.skt.Read(buf[read:])
		if err != nil {
			return -1, err
		}

		read += _read
	}

	return read, errors.New("failed to read from socket")
}
