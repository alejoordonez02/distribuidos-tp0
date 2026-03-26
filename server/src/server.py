from threading import Thread, Condition, Lock
import logging
import signal
from typing import Optional

from .net import Ack, Bet, Query, Response, Fin
from .has_won import has_won
from .net import Conn, Rendezvous
from .storage import load_bets, store_bets


class Server:
    def __init__(self, port, listen_backlog, agency_amount):
        self._keep_running = False
        self.listener = Rendezvous(("", port), listen_backlog)
        self.current: list[Conn] = []
        self.client_handlers: list[Thread] = []

        self.id_address_map: set[tuple[int, str]] = set()
        self.agency_amount = agency_amount
        self.pending = agency_amount
        self.results: dict[int, set[str]] = {}

        self.mtx = Lock()
        self.cv = Condition()

    def start(self):
        """
        Starts accepting incoming connections and handling client messages.
        """
        self._keep_running = True
        signal.signal(signal.SIGINT, self.stop)
        signal.signal(signal.SIGTERM, self.stop)
        self.listener.start()
        self.__run()

    def __run(self):
        while self._keep_running:
            if not (conn_info := self.listener.accept_connection()):
                continue

            conn, _ = conn_info
            self.current.append(conn)
            handler = Thread(target=self.__handle_client_connection, args=[conn])
            handler.start()
            self.client_handlers.append(handler)

    def __handle_client_connection(self, client: Conn):
        try:
            msg = client.recv()
        except Exception as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
            return

        if isinstance(msg, list):
            self.__handle_bets(msg, client)
        elif isinstance(msg, Fin):
            self.__handle_fin(msg, client)
        elif isinstance(msg, Query):
            self.__handle_query(msg, client)
        else:
            logging.error(
                f"action: receive_message | result: fail | error: unsupported message {msg.__dict__}"
            )

    def __handle_bets(self, bets: list[Bet], client: Conn):
        if not all(isinstance(b, Bet) and b.agency == bets[0].agency for b in bets):
            client.send(Ack(False))
            logging.error(
                "action: apuesta_recibida | result: fail | error: malformed message"
            )
            return

        client.send(Ack(True))

        client_id = bets[0].agency
        client_addr = client.peer_addr
        with self.mtx:
            self.__add_pending_client(client_id, client_addr)
            store_bets(bets)

        logging.info(
            f"action: apuesta_recibida | result: success | cantidad: {len(bets)}"
        )

    def __handle_fin(self, _: Fin, client: Conn):
        """
        Handles a `Fin` message from a client.

        Decreases the pending clients count by one, if it reaches zero then it
        computes the lottery results and notifies all threads waiting for the
        condition var.
        """
        client.send(Ack(True))
        with self.cv:
            self.pending -= 1
            if self.pending:
                return

            for b in load_bets():
                winners = self.results.get(b.agency, set(""))
                if has_won(b):
                    winners.add(b.document)

                self.results[b.agency] = winners

            self.cv.notify_all()

    def __handle_query(self, _: Query, client: Conn):
        with self.cv:
            self.cv.wait_for(lambda: not self.pending)

        client_addr = client.peer_addr
        client_id = self.__get_pending_client_id(client_addr)
        if not client_id:
            client.send(Ack(False))
            return

        client_results = self.results[client_id]
        response = Response(list(client_results))
        client.send(response)

    def __add_pending_client(self, client_id: int, addr: tuple[str, int]) -> None:
        """
        Adds a client to the list of clients to which results are sent at the end.

        As clients communicate with the server initiating different TCP sockets, ports
        may not be always the same, thus the IP address of the client is used for
        uniquely identifying it. The `client_id` is used for filtering the results
        that are to be sent to each client.
        """
        self.id_address_map.add((client_id, addr[0]))

    def __get_pending_client_id(self, client_addr: tuple[str, int]) -> Optional[int]:
        """
        Returns a registered client's id, that's associated with its agency number, if
        there is one.
        """
        for client_id, addr in self.id_address_map:
            if client_addr[0] == addr:
                return client_id

    def stop(self, _signum, _frame):
        """
        Stops the server by having it
        - stop accepting incoming connections,
        - stop its current connections, and
        - join its client handler threads.
        """
        logging.info("action: stop | result: in_progress")
        self._keep_running = False
        self.listener.stop()
        for c in self.current:
            c.close()

        for h in self.client_handlers:
            h.join()

        logging.info("action: stop | result: success")
