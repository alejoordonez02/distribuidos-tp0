package common

type Ack struct {
	Ok bool
}

func (a Ack) Type() byte {
	return MSG_ACK
}
