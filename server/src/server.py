import logging
import signal

from .bet import Bet
from .net import Conn, Rendezvous, SerialError
from .response import Response
from .storage import store_bets


class Server:
    def __init__(self, port, listen_backlog):
        self._keep_running = False
        self.listener = Rendezvous(("", port), listen_backlog)

    def start(self):
        self._keep_running = True
        signal.signal(signal.SIGINT, self.stop)
        signal.signal(signal.SIGTERM, self.stop)
        self.listener.start()
        self.__run()

    def __run(self):
        while self._keep_running:
            client = self.listener.accept_connection()
            if client:
                self.__handle_client_connection(client)

    def __handle_client_connection(self, client: Conn):
        try:
            msg = client.recv()
        except SerialError as e:
            logging.info(f"action: receive_message | result: fail | error: {e}")
            return

        if isinstance(msg, list) and all(isinstance(bet, Bet) for bet in msg):
            client.send(Response(True))
            store_bets(msg)
            logging.info(
                f"action: apuesta_recibida | result: success | cantidad: {len(msg)}"
            )

        else:
            raise RuntimeError(f"unsupported message {msg.__dict__}")

    def stop(self, _signum, _frame):
        self._keep_running = False
        logging.info("action: stop | result: in_progress")
        self.listener.stop()
        logging.info("action: stop | result: success")
