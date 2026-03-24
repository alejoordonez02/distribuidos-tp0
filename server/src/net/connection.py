from socket import socket
import logging


from .protocol import serialize, deserialize

LEN_SIZE = 2
# this byte order can be independent from the one in `protocol.py`
# since the preappending msg size protocol is on the conn protocol
# layer
BYTE_ORDER = "big"
BUF_SIZE = 1024


class Conn:
    def __init__(self, skt: socket, peer_addr: tuple[str, int]):
        self.skt = skt
        self.peer_addr = peer_addr

    def send(self, msg):
        serial = b""
        serial_msg = serialize(msg)
        len_msg = len(serial_msg).to_bytes(LEN_SIZE, byteorder=BYTE_ORDER)
        serial += len_msg
        serial += serial_msg

        self.skt.sendall(serial)

    def recv(self):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            serial = b""
            serial += self.__recv_exact(LEN_SIZE)
            len_msg = int.from_bytes(serial, byteorder=BYTE_ORDER)
            serial += self.__recv_exact(len_msg)
            msg = deserialize(serial[LEN_SIZE : LEN_SIZE + len_msg])
            logging.info(
                f"action: receive_message | result: success | ip: {self.peer_addr[0]}"
            )

            return msg
        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")

    def __recv_exact(self, amount: int) -> bytes:
        buf = b""
        missing = amount
        while missing:
            received = self.skt.recv(missing)
            buf += received
            missing -= len(received)

        return buf
