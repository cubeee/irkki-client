package irkki

import (
	"github.com/cubeee/irkki-client/event"
	"github.com/cubeee/irkki-client/irc"
)

func NewClient(cfg irc.Config) *irc.Client {
	client := &irc.Client{
		Config:   cfg,
		Handlers: make(map[string]map[int]func(irc.Connection, *event.Event)),
	}
	return client
}
