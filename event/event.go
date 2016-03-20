package event

import (
	"fmt"
)

const (
	CONNECTED         string = "004"
	ERR_NICKNAMEINUSE string = "433"
	DISCONNECTED      string = "DISCONNECTED"
	RAW_MESSAGE       string = "RAW_MESSAGE"
	PING              string = "PING"
	PONG              string = "PONG"
	JOIN              string = "JOIN"
	QUIT              string = "QUIT"
	NICK              string = "NICK"
	PASS              string = "PASS"
	USER              string = "USER"
)

type Event struct {
	Raw     string
	Source  string
	User    string
	Command string
	Args    []string
}

func (e Event) String() string {
	return fmt.Sprintf("\x1b[32;1mSource='%s', \n\t\x1b[33;1mUser='%s', \n\t\x1b[34;1mCommand='%s', \n\t\x1b[35;1mArgs='%s', \n\t\x1b[36;1mRaw='%s'\x1b[0m",
		e.Source, e.User, e.Command, e.Args, e.Raw)
}
