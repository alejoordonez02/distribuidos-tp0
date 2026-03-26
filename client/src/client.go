package src

import (
	"errors"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/op/go-logging"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/src/comms"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/src/comms/messages"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID             string
	ServerAddress  string
	LoopAmount     int
	LoopPeriod     time.Duration
	BatchMaxAmount int
}

// Client Entity that encapsulates how
type Client struct {
	config      ClientConfig
	conn        comms.Conn
	storage     *Storage
	keepRunning bool
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config:      config,
		keepRunning: false,
	}

	return client
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) Run() {
	err := c.createClientStorage()
	if err != nil {
		log.Criticalf(
			"action: create_storage | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}

	defer c.storage.Close()

	c.keepRunning = true
	go c.shouldKeepRunning()

	for {
		bets, err := c.storage.LoadBets(c.config.BatchMaxAmount, c.config.ID)
		if err == io.EOF && len(bets) == 0 {
			break
		}

		if err != io.EOF && err != nil {
			log.Errorf("action: load_bets | result: fail | client_id: %v | error: %v",
				c.config.ID, err)
			return
		}

		if err := c.sendBetAndRecvAck(bets); err != nil {
			log.Infof("action: send_bet | result: fail | client_id: %v | error: %v",
				c.config.ID, err)
			return
		}

		log.Infof("action: send_bet | result: success | client_id: %v", c.config.ID)
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)

	response, err := c.sendQueryAndReceiveResponse()
	if err != nil {
		log.Infof("action: consulta_ganadores | result: failed | client_id: %v | error: %v",
			c.config.ID, err)
		return
	}

	log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %v", response.WinnerAmount)
}

func (c *Client) sendBetAndRecvAck(bets messages.BetBatch) error {
	c.createClientConn()
	defer c.conn.Close()

	if err := c.conn.Send(bets); err != nil {
		return err
	}

	ack, err := c.conn.Recv()
	if err != nil {
		return err
	}

	switch m := ack.(type) {
	case messages.Ack:
		if !m.Ok {
			return errors.New("server responded with NACK")
		}
	default:
		return errors.New("unexpected message")
	}

	return nil
}

func (c *Client) sendQueryAndReceiveResponse() (*messages.Response, error) {
	c.createClientConn()
	defer c.conn.Close()

	query := messages.NewQuery()
	if err := c.conn.Send(query); err != nil {
		return nil, err
	}

	response, err := c.conn.Recv()
	if err != nil {
		return nil, err
	}

	switch m := response.(type) {
	case messages.Response:
		return &m, nil
	default:
		return nil, errors.New("unexpected message")
	}
}

func (c *Client) shouldKeepRunning() error {
	sigChan := make(chan os.Signal, 1)
	defer close(sigChan)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	log.Infof("action: stop | result: in_progress | client_id: %v",
		c.config.ID,
	)

	c.keepRunning = false
	err := c.conn.Close()
	if err != nil {
		return err
	}

	return nil
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientConn() error {
	conn, err := comms.NewConn(c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}

	c.conn = conn
	return nil
}

func (c *Client) createClientStorage() error {
	storage, err := NewStorage()
	if err != nil {
		return err
	}

	c.storage = storage
	return nil
}
