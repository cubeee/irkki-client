package event_test

import (
	"github.com/cubeee/irkki-client/event"
	"reflect"
	"testing"
)

func TestEventParse(t *testing.T) {
	tests := []struct {
		in  string
		out event.Event
	}{
		{":nick!ident@host.com PRIVMSG me :Hello", event.Event{":nick!ident@host.com PRIVMSG me :Hello", "nick!ident@host.com", "nick", "PRIVMSG", []string{"me", ":Hello"}}},
		{"PING :irc.server.com", event.Event{"PING :irc.server.com", "", "", "PING", []string{":irc.server.com"}}},
		{"PING", event.Event{"PING", "", "", "PING", []string{}}},
		{
			":irc.server.com 002 nick :Your host is irc.server.com[0.0.0.0/6667], running version ircd-seven-1.1.3",
			event.Event{
				":irc.server.com 002 nick :Your host is irc.server.com[0.0.0.0/6667], running version ircd-seven-1.1.3",
				"irc.server.com", "nick", "002", []string{":Your", "host", "is", "irc.server.com[0.0.0.0/6667],", "running", "version", "ircd-seven-1.1.3"}}},
	}

	for _, test := range tests {
		if actual, err := event.ParseEvent(test.in); err != nil {
			t.Fatal("Failed to parse event: %s", err.Error())
		} else {
			expected := test.out
			if actual.String() == "" {
				t.Error("Got an empty String(), non-empty expected")
			}
			if actual.Raw != expected.Raw {
				t.Errorf("Raw event mismatch, got '%s', expected '%s'", actual.Raw, expected.Raw)
			}
			if actual.Source != expected.Source {
				t.Errorf("Source does not match, got '%s', expected '%s'", actual.Source, expected.Source)
			}
			if actual.User != expected.User {
				t.Errorf("User does not match, got '%s', expected '%s'", actual.User, expected.User)
			}
			if actual.Command != expected.Command {
				t.Errorf("Command does not match, got '%s', expected '%s'", actual.Command, expected.Command)
			}
			if !reflect.DeepEqual(actual.Args, expected.Args) {
				t.Errorf("Args do not match, got '%s', expected '%s'", actual.Args, expected.Args)
			}
		}
	}
}
