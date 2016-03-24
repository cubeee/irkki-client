package irc

import (
	"bufio"
	"errors"
	"net"
	"sync"
)

type Connection struct {
	mutex  sync.RWMutex
	config *Config

	dialer *net.Dialer
	socket net.Conn
	io     *bufio.ReadWriter
	out    chan string

	connected bool
}

func (c *Connection) write(line string) error {
	if line == "" {
		return errors.New("Writing empty lines is not permitted")
	}
	if _, err := c.io.WriteString(line + "\r\n"); err != nil {
		return err
	}
	if err := c.io.Flush(); err != nil {
		return err
	}
	return nil
}

func NewConnection(cfg Config) *Connection {
	dialer := new(net.Dialer)
	dialer.Timeout = cfg.Timeout
	return &Connection{
		config: &cfg,
		dialer: dialer,
		out:    make(chan string, 64),
	}
}
