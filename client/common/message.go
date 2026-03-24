package common

const (
	MSG_RESPONSE = 0x10
	MSG_ACK      = 0x11
)

type Message interface {
	Type() byte
}
