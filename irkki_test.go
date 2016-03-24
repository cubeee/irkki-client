package irkki_test

import (
	"github.com/cubeee/irkki-client"
	"github.com/cubeee/irkki-client/irc"
	"testing"
)

func TestClientConstructor(t *testing.T) {
	cfg := irc.Config{
		Server: "irc.server.com",
	}
	client := irkki.NewClient(cfg)
	if client.Config.Server != "irc.server.com" {
		t.Fatal("Configs' Server does not match, got '%s', expected '%s'", client.Config.Server, "irc.server.com")
	}
}
