import logging
import signal

from .bet import Bet, has_won
from .net import Conn, Rendezvous, SerialError
from .ack import Ack
from .query import Query
from .storage import load_bets, store_bets
from .response import Response

AGENCIES = 5


class Server:
    def __init__(self, port, listen_backlog):
        self._keep_running = False
        self.listener = Rendezvous(("", port), listen_backlog)
        self.current: list[Conn] = []
        self.pending: list[tuple[int, tuple[str, int]]] = []

    def start(self):
        self._keep_running = True
        signal.signal(signal.SIGINT, self.stop)
        signal.signal(signal.SIGTERM, self.stop)
        self.listener.start()
        self.__run()

    def __add_pending(self, client_id: int, addr: tuple[str, int]) -> int | None:
        for c_id, a in self.pending:
            if (c_id, a) == (client_id, addr):
                return c_id

        self.pending.append((client_id, addr))

    def __get_client_id(self, client_addr: tuple[str, int]) -> int | None:
        for client_id, addr in self.pending:
            if client_addr == addr:
                return client_id

    def __run(self):
        pending_agencies = AGENCIES
        while self._keep_running and pending_agencies:
            if not (conn_info := self.listener.accept_connection()):
                continue

            conn, addr = conn_info
            self.current.append(conn)
            done = self.__handle_client_connection(conn, addr)
            pending_agencies -= done

        self.__send_results()

    def __handle_client_connection(self, client: Conn, addr: tuple[str, int]) -> bool:
        """
        Handles a client connection.

        Returns `True` if the client has completed its task.
        """
        try:
            msg = client.recv()
        except SerialError as e:
            logging.info(f"action: receive_message | result: fail | error: {e}")
            return False

        if isinstance(msg, list) and all(
            isinstance(bet, Bet) and bet.agency == msg[0].agency for bet in msg
        ):
            self.__add_pending(msg[0].agency, addr)
            client.send(Ack(True))
            store_bets(msg)
            logging.info(
                f"action: apuesta_recibida | result: success | cantidad: {len(msg)}"
            )
            return False
        elif isinstance(msg, Query):
            return True
        else:
            client.send(Ack(False))
            raise RuntimeError(f"unsupported message {msg.__dict__}")

    def __send_results(self):
        recipients = {}
        for b in load_bets():
            won = has_won(b)
            recipients[b.agency] += won

        for c in self.current:
            client_id = self.__get_client_id(c.peer_addr)
            winner_amount = recipients[client_id]
            response = Response(winner_amount)
            c.send(response)

    def stop(self, _signum, _frame):
        self._keep_running = False
        logging.info("action: stop | result: in_progress")
        self.listener.stop()
        for c in self.current:
            c.close()

        logging.info("action: stop | result: success")
