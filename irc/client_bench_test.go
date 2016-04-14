package irc

import (
	"fmt"
	"github.com/cubeee/irkki-client/event"
	"testing"
)

func BenchmarkClient(b *testing.B) {
	b.StopTimer()
	server := createServer()
	go server.Listen(6667, true, "bench client")

	client := createClient()
	client.Config.Server = "localhost"
	client.Config.Port = 6667

	client.HandleEvent(event.PING, func(conn Connection, evt *event.Event) {

	})

	if err := client.Connect(); err != nil {
		b.Fatalf("Failed to connect to server: %s", err.Error())
	}

	b.ResetTimer()
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		client.Conn.WriteRaw(fmt.Sprintf("PING :RAW %v", n))
	}
}
