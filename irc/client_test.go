package irc

import (
	"bufio"
	"crypto/tls"
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
	quit     chan bool
	conns    []net.Conn
	test     string
}

func createServer(test string) *MockServer {
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

	return &MockServer{
		text:  text,
		quit:  make(chan bool),
		conns: make([]net.Conn, 0, 10),
		test:  test,
	}
}

func (s *MockServer) Listen(port int, read bool, fn func(success bool, address string)) {
	defer func() {
		for _, conn := range s.conns {
			if conn != nil {
				conn.Close()
			}
		}
	}()

	listener, err := net.Listen("tcp4", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err.Error())
		fn(false, "")
		return
	}
	s.listener = listener
	fn(true, listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				return
			default:
			}
			continue
		}
		if conn == nil {
			continue
		}
		s.conns = append(s.conns, conn)
		if read {
			go s.handleRequest(conn, len(s.conns)-1)
		}
	}
}

func (s *MockServer) Close() {
	close(s.quit)
	if s.listener != nil {
		s.listener.Close()
	}
}

func (s *MockServer) handleRequest(conn net.Conn, id int) {
	defer func() {
		conn.Close()
		s.conns[id] = nil
	}()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			conn.Close()
			return
		}
		log.Println(s.test, "-->", message)
		if s.idx == 0 {
			for i := 0; i < len(s.text); i++ {
				conn.Write([]byte(s.text[i]))
			}
		}
		s.idx = s.idx + 1
	}
}

func createClient() *Client {
	cfg := *NewConfig(&User{"test", "test"})
	client := Client{
		Config:   cfg,
		Handlers: make(map[string]map[int]func(Connection, *event.Event)),
	}
	return &client
}

func createBasicConfiguredClient(server string, port int) *Client {
	client := createClient()
	client.Config.Server = server
	client.Config.Port = port
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

func TestClientConnectInvalidPort(t *testing.T) {
	client := createBasicConfiguredClient("localhost", 6667)
	defer client.Disconnect()

	ports := []int{-1, 0, 65536}
	for _, port := range ports {
		client.Config.Port = port
		if err := client.Connect(); err != nil {
			if strings.Index(err.Error(), "Invalid port given, 1-65535 expected, got") == -1 {
				t.Errorf("Client accepted invalid port %v", port)
			}
		}
	}
}

func TestClientConnect(t *testing.T) {
	server := createServer("connect")
	go server.Listen(0, true, func(success bool, address string) {
		defer server.Close()
		if !success {
			t.Fatal("Failed to bind server")
			return
		}

		client := createBasicConfiguredClient("localhost", 6667)
		defer client.Disconnect()

		if err := client.Connect(); err != nil {
			t.Fatalf("Client failed to connect to server: %s", err.Error())
		}

		if !client.Connected() {
			t.Error("Client reports not being connected")
		}
	})
}

func TestClientConnectWithAddress(t *testing.T) {
	cases := []struct{ Address, Error string }{
		{"", "Address of following format required: <server>:<port>"},
		{":", "Failed to parse a numeric port from given '', make sure to use a numeric port and the following address format: <server>:<port>"},
		{"localhost:", "Failed to parse a numeric port from given '', make sure to use a numeric port and the following address format: <server>:<port>"},
		{":6667", "Empty server address given"},
		{"localhost:6667", ""},
	}

	for _, c := range cases {
		server := createServer("connect_with_address")
		go server.Listen(0, true, func(success bool, address string) {
			defer server.Close()
			if !success {
				t.Fatal("Failed to bind server")
				return
			}

			client := createClient()
			if err := client.ConnectWithAddress(c.Address); err != nil {
				if err.Error() != c.Error {
					t.Errorf("Testing address '%s' failed: %s", c.Address, err.Error())
				}
			}
		})
	}
}

func TestClientDisconnect(t *testing.T) {
	server := createServer("disconnect")
	go server.Listen(0, true, func(success bool, address string) {
		defer server.Close()
		if !success {
			t.Fatal("Failed to bind server")
			return
		}

		client := createClient()
		defer client.Disconnect()

		if err := client.ConnectWithAddress(address); err != nil {
			t.Fatalf("Client failed to connect to server: %s", err.Error())
		}

		if err := client.Disconnect(); err != nil {
			t.Error("Failed to disconnect")
		}

		if client.Connected() {
			t.Error("Client reports being connected after disconnecting")
		}
	})
}

func TestClientNoSSLConfig(t *testing.T) {
	server := createServer("no_ssl")
	go server.Listen(0, true, func(success bool, address string) {
		defer server.Close()
		if !success {
			t.Fatal("Failed to bind server")
			return
		}

		client := createClient()
		client.Config.SSL = true
		defer client.Disconnect()

		if err := client.ConnectWithAddress(address); err == nil {
			t.Fatal("Client did not complain about missing SSL config when expected to")
		}
	})
}

func TestClientSSL(t *testing.T) {
	server := createServer("ssl")
	go server.Listen(0, true, func(success bool, address string) {
		defer server.Close()
		if !success {
			t.Fatal("Failed to bind server")
			return
		}

		client := createClient()
		client.Config.SSL = true
		client.Config.SSLConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
		defer client.Disconnect()

		if err := client.ConnectWithAddress(address); err != nil {
			t.Fatalf("Client failed to connect to server: %s", err.Error())
		}

		if !client.Connected() {
			t.Error("Client is not connected")
		}
	})
}
