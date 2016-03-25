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
	Handlers  map[string]map[int]func(Connection, *event.Event)
}

type EventHandler struct {
	id int64
	fn func(Connection, *event.Event)
}

type EventHandlers struct {
	id int64
	Handlers map[string][]EventHandler
	sync.RWMutex
}

func (c Client) HandleEvent(evt string, fn func(Connection, *event.Event)) int {
	evt = strings.ToUpper(evt)
	id := 0
	if _, ok := c.Handlers[evt]; !ok {
		c.Handlers[evt] = make(map[int]func(Connection, *event.Event))
		id = 0
	} else {
		id = len(c.Handlers[evt])
	}
	c.Handlers[evt][id] = fn
	return id
}

func (c Client) RemoveEventHandler(evt string, id int) bool {
	evt = strings.ToUpper(evt)
	if e, ok := c.Handlers[evt];  ok {
		if _, ok := e[id]; ok {
			delete(c.Handlers[evt], id)
			return true
		}
		return false
	}
	return false
}

func (c Client) fireEvent(evt *event.Event) {
	if evt.Command != event.RAW_MESSAGE {
		fmt.Println(evt)
	}
	command := strings.ToUpper(evt.Command)
	if handlers, ok :=  c.Handlers[command]; ok {
		for _, handler := range handlers {
			handler(c.Conn, evt)
		}
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

			if evt, err := event.ParseEvent(line); err == nil {
				if evt.Command == event.PING {
					source := strings.Join(evt.Args[1:], " ")
					c.Conn.Pong(source)
				}
				c.fireEvent(evt)
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
				log.Panicln("Failed to send!")
				return
			}
		}
	}
}
