package irc

import (
	"bufio"
	"errors"
	"golang.org/x/net/proxy"
	"net"
	"sync"
)

type Connection struct {
	mutex  sync.RWMutex
	config *Config

	dialer      *net.Dialer
	proxyDialer proxy.Dialer
	socket      net.Conn
	io          *bufio.ReadWriter
	out         chan string

	connected bool
}

func (c *Connection) write(line string) error {
	if line == "" {
		return errors.New("Writing empty lines is not permitted")
	}
	c.io.WriteString(line + "\r\n")
	return c.io.Flush()
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
