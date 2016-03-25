package irc

import (
	"github.com/cubeee/irkki-client/event"
	"testing"
	"time"
)

func createClient() Client {
	cfg := *NewConfig(&User{"test", "test"})
	client := Client{
		Config:   cfg,
		Handlers: make(map[string]map[int]func(Connection, *event.Event)),
	}
	return client
}

func TestClientEventHandler(t *testing.T) {
	done := make(chan bool)

	client := createClient()
	client.HandleEvent(event.CONNECTED, func(conn Connection, event *event.Event) {
		<-done
	})
	select {
	case <-done:
		client.fireEvent(&event.Event{
			Command: event.CONNECTED,
		})
		t.Log("hello")
	case <-time.After(100 * time.Millisecond):
	case <-time.After(2 * time.Second):
		t.Error("No reply from async event handler after 2 second grace period")
	}
}