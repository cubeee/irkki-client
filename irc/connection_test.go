package irc

import (
	"bufio"
	"bytes"
	"testing"
)

func createRWConnection() (*Connection, *bytes.Buffer) {
	var buf bytes.Buffer
	con := createConnection()
	reader := bufio.NewReader(con.socket)
	writer := bufio.NewWriter(&buf)
	con.io = bufio.NewReadWriter(reader, writer)
	return con, &buf
}

func TestConnectionEmptyWrite(t *testing.T) {
	con, _ := createRWConnection()

	if err := con.write(""); err == nil {
		t.Error("Expected error from writing empty string to connection")
	}
}

func TestConnectionWrite(t *testing.T) {
	con, buf := createRWConnection()

	line := "test non-empty write"
	expected := line + "\r\n"
	if err := con.write(line); err != nil {
		t.Error("Expected no error from writing non-empty string to connection")
	} else {
		str := buf.String()
		if str != expected {
			t.Errorf("Connection read mismatch, got '%s', expected '%s'", str, expected)
		}
	}
}
