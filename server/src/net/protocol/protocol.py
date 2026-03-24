from enum import Enum

BYTE_ORDER = "big"

LEN_TYPE = 1
LEN_STR_SIZE = 1  # preappended string length comes in one byte

LEN_BATCH_SIZE = 2
LEN_BET_NUMBER = 8


class MsgType(Enum):
    TYPE_BET_BATCH = b"\x00"
    TYPE_RESPONSE_ACK = b"\x01"
    TYPE_RESPONSE_NACK = b"\x02"
