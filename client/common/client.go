package common

import (
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/op/go-logging"
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
	conn        Conn
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
	c.createClientSocket()
	c.createClientStorage()
	defer c.conn.Close()
	defer c.storage.Close()
	c.keepRunning = true
	go c.shouldKeepRunning()

	time.Sleep(time.Second * 10)
	bets, err := c.storage.LoadBets(c.config.BatchMaxAmount)
	if err != io.EOF && err != nil {
		log.Errorf("action: load_bets | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}

	if err := c.conn.Send(bets); err != nil {
		log.Errorf("action: send_message | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}

	response, err := c.conn.Recv()
	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}

	if response.Ack {
		// log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v",
		// 	c.bet.Document, c.bet.Number,
		// )
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
func (c *Client) createClientSocket() error {
	conn, err := NewConn(c.config.ServerAddress)
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
	// las variables constanteadas de storage deberían venir
	// por cfg pero fiaca ahora mismo
	storage, err := NewStorage()
	if err != nil {
		log.Criticalf(
			"action: create_storage | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}

	c.storage = storage
	return nil
}
