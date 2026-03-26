package protocol

const (
	TYPE_BET_BATCH = 0x00
	TYPE_ACK       = 0x01
	TYPE_NACK      = 0x02
	TYPE_QUERY     = 0x03
	TYPE_RESPONSE  = 0x04
	TYPE_FIN       = 0x05

	LEN_WINNER_AMOUNT = 2
	LEN_STR_SIZE      = 1
)
