import logging
import signal

from .bet import Bet
from .net import Conn, Rendezvous
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
        msg = client.recv()
        if isinstance(msg, Bet):
            client.send(Response(True))
            store_bets([msg])
            logging.info(
                f"action: apuesta_almacenada | result: success | dni: {msg.document} | numero: {msg.number}"
            )

        else:
            raise RuntimeError(f"unsupported message {msg.__dict__}")

    def stop(self, _signum, _frame):
        self._keep_running = False
        logging.info("action: stop | result: in_progress")
        self.listener.stop()
        logging.info("action: stop | result: success")
