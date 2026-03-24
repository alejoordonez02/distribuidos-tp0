from enum import Enum

BYTE_ORDER = "big"

LEN_TYPE = 1
LEN_STR_SIZE = 1  # preappended string length comes in one byte

LEN_BATCH_SIZE = 2
LEN_BET_NUMBER = 8
LEN_WINNER_AMOUNT = 2


class MsgType(Enum):
    TYPE_BET_BATCH = b"\x00"
    TYPE_ACK = b"\x01"
    TYPE_NACK = b"\x02"
    TYPE_QUERY = b"\x03"
    TYPE_RESPONSE = b"\x04"
