from socket import socket, AF_INET, SOCK_STREAM, SHUT_RDWR
import logging

from .connection import Conn


class Rendezvous:
    def __init__(self, addr: tuple[str, int], listen_backlog: int):
        self._keep_running = False
        self.skt = socket(AF_INET, SOCK_STREAM)
        self.skt.bind(addr)
        self.listen_backlog = listen_backlog

    def start(self):
        self.skt.listen(self.listen_backlog)

    def accept_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """
        try:
            logging.info("action: accept_connections | result: in_progress")
            c, addr = self.skt.accept()
            logging.info(
                f"action: accept_connections | result: success | ip: {addr[0]}"
            )
            return Conn(c, addr)
        except OSError as e:
            if self._keep_running:
                logging.info(f"action: accept_connections | result: fail | error: {e}")

    def stop(self):
        self._keep_running = False
        self.skt.shutdown(SHUT_RDWR)
        self.skt.close()
