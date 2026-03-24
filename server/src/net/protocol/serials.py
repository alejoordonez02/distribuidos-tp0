from ...response import Response
from ...bet import Bet
from .protocol import BYTE_ORDER, LEN_STR_SIZE, MsgType, LEN_TYPE


def serialize(msg) -> bytes:
    if isinstance(msg, Response):
        return __serialize_response(msg)
    else:
        raise RuntimeError(f"unsupported serialization for message {msg}")


def deserialize(serial: bytes):
    msg_type = serial[:LEN_TYPE]
    if msg_type == MsgType.TYPE_BET.value:
        return __deserialize_bet(serial[LEN_TYPE:])
    else:
        raise RuntimeError(f"unknown message type {msg_type}")


def __serialize_response(response: Response) -> bytes:
    if response.ack:
        serial = MsgType.TYPE_RESPONSE_ACK.value
    else:
        serial = MsgType.TYPE_RESPONSE_NACK.value

    return serial


def __deserialize_bet(serial: bytes) -> Bet:
    ptr = 0
    agency, consumed = __deserialize_string(serial[ptr:])
    ptr += consumed
    first_name, consumed = __deserialize_string(serial[ptr:])
    ptr += consumed
    last_name, consumed = __deserialize_string(serial[ptr:])
    ptr += consumed
    document, consumed = __deserialize_string(serial[ptr:])
    ptr += consumed
    birthdate, consumed = __deserialize_string(serial[ptr:])
    ptr += consumed
    number, consumed = __deserialize_string(serial[ptr:])

    bet = Bet(agency, first_name, last_name, document, birthdate, number)
    return bet


def __deserialize_string(serial: bytes) -> tuple[str, int]:
    """
    Deserializes a byte slice into a string.

    Returns the deserialized string along with the bytes it consumed,
    including the ones for its length.
    """
    length = int.from_bytes(serial[:LEN_STR_SIZE], byteorder=BYTE_ORDER)
    len_total = LEN_STR_SIZE + length
    string = serial[LEN_STR_SIZE:len_total].decode()

    return string, len_total
