package irc

import (
	"testing"
)

func createConnection() *Connection {
	cfg := *NewConfig(&User{"test", "test"})
	return NewConnection(cfg)
}

func testCommand(t *testing.T, expected, functionName string, fn func(con *Connection, args ...interface{}), args ...interface{}) {
	con := createConnection()
	fn(con, args...)
	select {
	case msg, ok := <-con.out:
		if ok {
			if msg != expected {
				t.Errorf("%s expected '%s', got '%s'", functionName, expected, msg)
			}
		} else {
			t.Errorf("Failed to read from channel for '%s'", functionName)
		}
	default:
		if expected != "" {
			t.Errorf("Nothing received from channel for '%s'", functionName)
		}
	}
}

func TestRaw(t *testing.T) {
	expected := "TEST WRITE_RAW"
	testCommand(t, expected, "WriteRaw", func(con *Connection, args ...interface{}) {
		con.WriteRaw(args[0].(string))
	}, expected)
}

func TestPass(t *testing.T) {
	testCommand(t, "PASS password123!test", "Pass", func(con *Connection, args ...interface{}) {
		con.Pass(args[0].(string))
	}, "password123!test")
}

func TestPong(t *testing.T) {
	testCommand(t, "PONG :irc.server.com", "Pong", func(con *Connection, args ...interface{}) {
		con.Pong(args[0].(string))
	}, ":irc.server.com")
}

func TestUser(t *testing.T) {
	testCommand(t, "USER user_test 8 * :full_name_test", "User", func(con *Connection, args ...interface{}) {
		con.User(args[0].(string), args[1].(int), args[2].(string))
	}, "user_test", 8, "full_name_test")
}

func TestNick(t *testing.T) {
	testCommand(t, "NICK user_test", "Nick", func(con *Connection, args ...interface{}) {
		con.Nick(args[0].(string))
	}, "user_test")
}

func TestJoin(t *testing.T) {
	testCommand(t, "JOIN #channel", "Join", func(con *Connection, args ...interface{}) {
		con.Join(args[0].(string))
	}, "#channel")
}

func TestQuit(t *testing.T) {
	testCommand(t, "QUIT :bye", "Quit", func(con *Connection, args ...interface{}) {
		con.Quit(args[0].(string))
	}, "bye")
}

func TestMessage(t *testing.T) {
	testCommand(t, "PRIVMSG #test :Test", "Message", func(con *Connection, args ...interface{}) {
		con.Message(args[0].(string), args[1].(string))
	}, "#test", "Test")
}

func TestMessagef(t *testing.T) {
	testCommand(t, "PRIVMSG #test :Test: Success", "Messagef", func(con *Connection, args ...interface{}) {
		con.Messagef(args[0].(string), args[1].(string), args[2:]...)
	}, "#test", "%s: %s", "Test", "Success")
}