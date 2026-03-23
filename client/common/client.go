package common

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	BetNumber     uint64
	BetName       string
	BetSurname    string
	BetBirth      string
}

// Client Entity that encapsulates how
type Client struct {
	config      ClientConfig
	bet         Bet
	conn        Conn
	keepRunning bool
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	bet := NewBet(config.BetNumber, config.BetName, config.BetSurname, config.BetBirth)

	client := &Client{
		config:      config,
		bet:         bet,
		keepRunning: false,
	}
	return client
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

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) Run() {
	c.createClientSocket()
	defer c.conn.Close()
	c.keepRunning = true
	go c.shouldKeepRunning()

	err := c.conn.Send(&c.bet)
	if err != nil {
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
		log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v",
			0, 0, // TODO
		)
	}
}
