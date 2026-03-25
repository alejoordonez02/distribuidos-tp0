import threading
import logging
import signal
from typing import Optional

from .net import Ack, Bet, Query, Response
from .has_won import has_won
from .net import Conn, Rendezvous, SerialError
from .storage import load_bets, store_bets


class Server:
    def __init__(self, port, listen_backlog, agency_amount):
        self._keep_running = False
        self.listener = Rendezvous(("", port), listen_backlog)
        self.current: list[Conn] = []
        self.pending: set[tuple[int, str]] = set()  # agency, ip
        self.agency_amount = agency_amount
        self.mtx = threading.Lock()
        self.sem = threading.Semaphore(0)

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
        wait_to_send_results = threading.Thread(target=self.__send_results)
        wait_to_send_results.start()

        while self._keep_running:
            if not (conn_info := self.listener.accept_connection()):
                continue

            conn, addr = conn_info
            self.current.append(conn)
            client_handle = threading.Thread(
                target=self.__handle_client_connection, args=(conn, addr)
            )
            client_handle.start()

    def __handle_client_connection(self, client: Conn, addr: tuple[str, int]):
        """
        Handles a client connection.
        """
        try:
            msg = client.recv()
        except SerialError as e:
            logging.info(f"action: receive_message | result: fail | error: {e}")
            return

        if isinstance(msg, list) and all(
            isinstance(bet, Bet) and bet.agency == msg[0].agency for bet in msg
        ):
            self.__handle_bets(msg, client, addr)
        elif isinstance(msg, Query):
            self.__handle_query(msg)
        else:
            client.send(Ack(False))
            raise RuntimeError(f"unsupported message {msg.__dict__}")

    def __handle_bets(self, bets: list[Bet], client: Conn, addr: tuple[str, int]):
        agency_number = bets[0].agency
        self.__add_pending(agency_number, addr)
        client.send(Ack(True))
        with self.mtx:
            store_bets(bets)

        logging.info(
            f"action: apuesta_recibida | result: success | cantidad: {len(bets)}"
        )

    def __handle_query(self, _: Query):
        self.sem.release()

    def __send_results(self):
        """
        Sends each client their corresponding lottery results, once all of them are
        done sending bets, after closing the `listener`.
        """
        for _ in range(0, self.agency_amount):
            self.sem.acquire()

        self.__stop_listening()

        recipients = {}
        for b in load_bets():
            won = has_won(b)
            recipients[b.agency] = recipients.get(b.agency, 0) + won

        for c in self.current:
            client_id = self.__get_client_id(c.peer_addr)
            winner_amount = recipients[client_id]
            response = Response(winner_amount)
            c.send(response)

    def __stop_listening(self):
        self._keep_running = False
        self.listener.stop()

    def stop(self, _signum, _frame):
        logging.info("action: stop | result: in_progress")
        self.__stop_listening()
        for c in self.current:
            c.close()

        logging.info("action: stop | result: success")
