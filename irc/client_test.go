package irc

import (
	"fmt"
	"github.com/cubeee/irkki-client/event"
	"github.com/cubeee/irkki-client/log"
	"net"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

type MockServer struct {
	listener net.Listener
	idx      int
	text     []string
}

func createServer() MockServer {
	text := []string{
		":irc.mockserver.com NOTICE * :*** Looking up your hostname...",
		":irc.mockserver.com NOTICE * :*** Checking Ident",
		":irc.mockserver.com NOTICE * :*** No Ident response",
		":irc.mockserver.com NOTICE * :*** Found your hostname",
		":irc.mockserver.com 001 nickname :Welcome to the freenode Internet Relay Chat Network nickname",
		":irc.mockserver.com 002 nickname :Your host is irc.mockserver.com[0.0.0.0/6667], running version mock-server-0.1",
		":irc.mockserver.com 003 nickname :This server was created Mon Jan 1 2000 at 00:00:00 UTC",
		":irc.mockserver.com 004 nickname irc.mockserver.com mock-server-0.1",
		":irc.mockserver.com 251 nickname :There are 0 users and 0 invisible on 1 servers",
		":irc.mockserver.com 252 nickname 0 :IRC Operators online",
	}
	return MockServer{
		text: text,
	}
}

func (s MockServer) Listen(port int, read bool) error {
	listener, err := net.Listen("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		return err
	}
	defer listener.Close()

	s.listener = listener
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("accept failed bruh: %s\n", err.Error())
			s.Close()
			break
		}
		if read {
			s.handleRequest(conn)
		}
	}
	return nil
}

func (s MockServer) handleRequest(conn net.Conn) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		log.Println("cant read shit", err.Error())
	}
	log.Println(s.idx, string(buf))

	s.idx = s.idx + 1
}

func (s MockServer) Close() error {
	if s.listener == nil {
		return nil
	}
	return s.listener.Close()
}

func createClient() Client {
	cfg := *NewConfig(&User{"test", "test"})
	client := Client{
		Config:   cfg,
		Handlers: make(map[string]map[int]func(Connection, *event.Event)),
	}
	return client
}

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false
	case <-time.After(timeout):
		return true
	}
}

func TestClientEventHandler(t *testing.T) {
	done := make(chan int, 5)
	var wg sync.WaitGroup
	wg.Add(1)

	client := createClient()
	client.HandleEvent(event.CONNECTED, func(conn Connection, event *event.Event) {
		done <- 1
		wg.Done()
	})
	client.fireEvent(&event.Event{
		Command: event.CONNECTED,
	})

	if waitTimeout(&wg, time.Second) {
		t.Error("No reply received from an event handler")
	} else {
		expected := 1
		actual := <-done
		if actual != expected {
			t.Errorf("Wrong reply received from event handler, got '%v', expected '%v'", actual, expected)
		}
	}
}

func TestMultipleRegisteredEventHandlers(t *testing.T) {
	client := createClient()
	firstEvent := client.HandleEvent(event.CONNECTED, func(conn Connection, event *event.Event) {})
	secondEvent := client.HandleEvent(event.CONNECTED, func(conn Connection, event *event.Event) {})
	expected := 0
	if firstEvent != expected {
		t.Errorf("Wrong event handler id received from first event, got '%v', expected '%v'", firstEvent, expected)
	}
	expected = 1
	if secondEvent != expected {
		t.Errorf("Wrong event handler id received from second event, got '%v', expected '%v'", secondEvent, expected)
	}
}

func TestRemovedEventHandler(t *testing.T) {
	done := make(chan int, 5)
	client := createClient()
	done <- 1
	handlerId := client.HandleEvent(event.CONNECTED, func(conn Connection, event *event.Event) {
		done <- 2
	})
	if !client.RemoveEventHandler(event.CONNECTED, handlerId) {
		t.Error("Failed to remove event handler")
	}
	client.fireEvent(&event.Event{
		Command: event.CONNECTED,
	})
	if len(done) != 1 {
		t.Error("Received an unexpected response from a supposedly removed event handler")
	}
}

func TestRemoveInvalidEventHandler(t *testing.T) {
	client := createClient()
	if client.RemoveEventHandler("INVALID_EVENT", 0) {
		t.Error("Removed non-existing event handler when it shouldn't have")
	}

	handlerId := client.HandleEvent(event.CONNECTED, func(conn Connection, event *event.Event) {

	})
	if client.RemoveEventHandler(event.CONNECTED, handlerId+1) {
		t.Error("Removed non-existing event handler with valid event name but invalid handler id")
	}
}

func TestClientConnectEmptyServerAddress(t *testing.T) {
	client := createClient()
	client.Config.Server = ""
	if err := client.Connect(); err == nil {
		t.Fatal("Client accepted empty server address")
	}
}

func TestClientConnectValidPort(t *testing.T) {
	server := createServer()
	go server.Listen(6667, true)

	client := createClient()
	client.Config.Server = "localhost"

	ports := []int{1, 65535}
	for _, port := range ports {
		client.Config.Port = port
		if err := client.Connect(); err != nil {
			if strings.Index(err.Error(), "connection refused") == -1 {
				t.Errorf("Client didn't accept valid port %v or returned another error: %s", port, err.Error())
			}
		}
		client.Disconnect()
	}
	server.Close()
}

func TestClientConnectInvalidPort(t *testing.T) {
	server := createServer()
	go server.Listen(6667, true)

	client := createClient()
	client.Config.Server = "localhost"
	client.Config.Port = 6667

	ports := []int{-1, 0, 65536}
	for _, port := range ports {
		client.Config.Port = port
		if err := client.Connect(); err != nil {
			if strings.Index(err.Error(), "Invalid port given, 1-65535 expected, got") == -1 {
				t.Errorf("Client accepted invalid port %v", port)
			}
		}
	}
	server.Close()
}
