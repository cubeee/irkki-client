package irkki

import (
	"github.com/cubeee/irkki-client/irc"
)

func NewClient(cfg irc.Config) irc.Client {
	client := irc.Client{
		Config:   cfg,
		Handlers: new(irc.EventHandlers),
	}
	client.Handlers.Handlers = make(map[string][]irc.EventHandler)
	return client
}
