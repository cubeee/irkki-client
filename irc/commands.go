package irc

import (
	"fmt"
	"github.com/cubeee/irkki-client/event"
)

func (c *Connection) WriteRaw(line string) {
	c.out <- line
}

func (c *Connection) Pass(password string) {
	c.WriteRaw(event.PASS + " " + password)
}

func (c *Connection) Pong(source string) {
	c.WriteRaw(event.PONG + " " + source)
}

func (c *Connection) User(name string, mode int, fullname string) {
	c.WriteRaw(fmt.Sprintf("%s %s %d * :%s", event.USER, name, mode, fullname))
}

func (c *Connection) Nick(nickname string) {
	c.WriteRaw(event.NICK + " " + nickname)
}

func (c *Connection) Join(channel string) {
	c.WriteRaw(event.JOIN + " " + channel)
}

func (c *Connection) Quit(message string) {
	c.WriteRaw(event.QUIT + " :" + message)
}
