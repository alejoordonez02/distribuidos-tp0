import logging
import signal
from typing import Optional

from .ack import Ack
from .bet import Bet, has_won
from .net import Conn, Rendezvous, SerialError
from .query import Query
from .response import Response
from .storage import load_bets, store_bets


class Server:
    def __init__(self, port, listen_backlog, agency_amount):
        self._keep_running = False
        self.listener = Rendezvous(("", port), listen_backlog)
        self.current: list[Conn] = []
        self.pending: set[tuple[int, str]] = set()  # agency, ip
        self.agency_amount = agency_amount

    def start(self):
        self._keep_running = True
        signal.signal(signal.SIGINT, self.stop)
        signal.signal(signal.SIGTERM, self.stop)
        self.listener.start()
        self.__run()

    def __add_pending(self, client_id: int, addr: tuple[str, int]) -> None:
        """
        Adds a client to the list of clients to which results are sent at the end.

        As clients communicate with the server initiating different TCP sockets, ports
        may not be always the same, thus the IP address of the client is used for
        uniquely identifying it. The `client_id` is used for filtering the results that
        are to be sent to each client.
        """
        self.pending.add((client_id, addr[0]))

    def __get_client_id(self, client_addr: tuple[str, int]) -> Optional[int]:
        for client_id, addr in self.pending:
            if client_addr[0] == addr:
                return client_id

    def __run(self):
        pending_agencies = self.agency_amount
        while self._keep_running and pending_agencies:
            if not (conn_info := self.listener.accept_connection()):
                continue

            conn, addr = conn_info
            self.current.append(conn)
            done = self.__handle_client_connection(conn, addr)
            pending_agencies -= done

        if self._keep_running:
            logging.info("action: sorteo | result: success")
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
            recipients[b.agency] = recipients.get(b.agency, 0) + won

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
