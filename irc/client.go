package irc

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"github.com/cubeee/irkki-client/event"
	"github.com/cubeee/irkki-client/log"
	"net"
	"strconv"
	"strings"
	"sync"
)

type Client struct {
	Conn      Connection
	Config    Config
	connected bool
	Handlers  *CommandHandlers
}

type EventHandler func(Connection, *event.Event)

type CommandHandlers struct {
	Handlers map[string][]EventHandler
	sync.RWMutex
}

func (c Client) HandleCommand(event string, handler EventHandler) {
	c.Handlers.Lock()
	c.Handlers.Handlers[event] = append(c.Handlers.Handlers[event], handler)
	c.Handlers.Unlock()
}

func (c Client) fireEvent(evt *event.Event) {
	if evt.Command != event.RAW_MESSAGE {
		fmt.Println(evt)
	}
	c.Handlers.RLock()
	defer c.Handlers.RUnlock()
	handlers, _ := c.Handlers.Handlers[evt.Command]
	for _, handler := range handlers {
		handler(c.Conn, evt)
	}
}

func (c Client) Connect() error {
	connection := *NewConnection(c.Config)
	c.Conn = connection

	server := net.JoinHostPort(c.Config.Server, strconv.Itoa(c.Config.Port))
	log.Println("Connecting to", server)
	// todo: check server address
	// todo: proxy
	if socket, err := c.Conn.dialer.Dial("tcp", server); err != nil {
		return err
	} else {
		c.Conn.socket = socket
		if c.Config.SSL {
			c.Conn.socket = tls.Client(c.Conn.socket, c.Config.SSLConfig)
		}
		c.postConnect(socket)

		if c.Config.Password != "" {
			c.Conn.Pass(c.Config.Password)
		}
		// todo: mask
		c.Conn.User(c.Config.User.Username, 0, c.Config.User.Realname)
		c.Conn.Nick(c.Config.User.Username)
	}
	return nil
}

func (c Client) postConnect(socket net.Conn) {
	c.Conn.io = bufio.NewReadWriter(
		bufio.NewReader(socket),
		bufio.NewWriter(socket))
	c.connected = true
	go c.send()
	go c.receive()
}

func (c Client) parseMessage(raw string) (*event.Event, error) {
	evt := &event.Event{Raw: raw}

	var command string
	var args []string

	if raw[0] == ':' {
		parts := strings.SplitN(raw[1:], " ", 2)
		source := parts[0]
		evt.Source = source
		if strings.Index(source, "!") == -1 {
			// from server
		} else {
			// user's
		}
		parts = strings.Split(parts[1], " ")
		command = parts[0]
		args = parts[1:]
	} else { // command
		parts := strings.Split(raw, " ")
		command = parts[0]
		args = parts[1:]

		if command == event.PING {
			source := strings.Join(parts[1:], " ")
			c.Conn.Pong(source)
			return evt, nil
		}
	}
	evt.Command = command
	evt.Args = args
	return evt, nil
}

func (c Client) receive() {
	disconnectEvent := &event.Event{
		Command: event.DISCONNECTED,
	}
	rawMessageEvent := &event.Event{
		Command: event.RAW_MESSAGE,
	}
	for c.connected {
		// todo: read timeout, socket.SetReadDeadline
		if line, err := c.Conn.io.ReadString('\n'); err != nil {
			disconnectEvent.Source = c.Config.Server
			c.fireEvent(disconnectEvent)
			// do we want to reconnect here or have events do it?
			// maybe have a flag in config for it
			// at least put some threshold here rofl
			// log.Println("Lost connection, reconnecting...")
			// c.Connect()
		} else {
			line = strings.Trim(line, "\r\n")
			rawMessageEvent.Raw = line
			c.fireEvent(rawMessageEvent)

			if event, err := c.parseMessage(line); err == nil {
				c.fireEvent(event)
			} else {
				log.Panicln("shit cant parse ", line)
			}
		}
	}
}

func (c Client) send() {
	for {
		select {
		case line := <-c.Conn.out:
			if err := c.Conn.write(line); err != nil {
				log.Panicln("Failed to send!!!")
				return
			}
		}
	}
}
