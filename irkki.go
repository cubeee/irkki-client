package irkki

import (
	"github.com/cubeee/irkki-client/irc"
	"github.com/cubeee/irkki-client/event"
)

func NewClient(cfg irc.Config) irc.Client {
	client := irc.Client{
		Config:   cfg,
		Handlers: make(map[string]map[int]func(irc.Connection, *event.Event)),
	}
	return client
}
