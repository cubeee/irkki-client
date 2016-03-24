# irkki-client

[![Travis Build status][travis-build-status-img]][travis-build-status] [![Drone Build status][drone-build-status-img]][drone-build-status] [![Coverage Status][coveralls-badge-img]][coveralls-url]

Irkki-client is a ridiculously named IRC client library written in go.
This library was developed primarily to be the underlaying IRC layer for a bouncer.

___This is still heavily under development and should not be used in production!___

## Usage
    go get github.com/cubeee/irkki-client

---

    import (
        "github.com/cubeee/irkki-client"
        "github.com/cubeee/irkki-client/event"
        "github.com/cubeee/irkki-client/irc"
        "github.com/cubeee/irkki-client/log"
        "time"
    )

    func main() {
        user := &irc.User{
            Username: "irkki-client",
            Realname: "irkki-client",
        }
        cfg := *irc.NewConfig(user)
        cfg.Server = "irc.freenode.net"
        cfg.Port = 6667
        cfg.Timeout = 60 * time.Second

        client := irkki.NewClient(cfg)
        client.HandleCommand(event.CONNECTED, func(conn irc.Connection, event *event.Event) {
            conn.Join("#channel")
        })
        err := client.Connect()
        if err != nil {
            log.Panicln("Failed to connect!")
        }
    }

[travis-build-status-img]: https://travis-ci.org/cubeee/irkki-client.svg
[travis-build-status]: https://travis-ci.org/cubeee/irkki-client
[drone-build-status-img]: https://drone.io/github.com/cubeee/irkki-client/status.png
[drone-build-status]: https://drone.io/github.com/cubeee/irkki-client/latest
[coveralls-badge-img]: https://coveralls.io/repos/github/cubeee/irkki-client/badge.svg?branch=master
[coveralls-url]: https://coveralls.io/github/cubeee/irkki-client?branch=master
