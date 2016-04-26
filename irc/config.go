package irc

import (
	"crypto/tls"
	"time"
)

type User struct {
	Username string
	Realname string
}

type Config struct {
	User *User

	Server       string
	Port         int
	Password     string
	SSL          bool
	SSLConfig    *tls.Config
	Proxy        bool
	ProxyAddress string
	Timeout      time.Duration
	QuitMessage  string
}

func NewConfig(user *User) *Config {
	cfg := &Config{
		User:        user,
		QuitMessage: "Bye!",
	}
	return cfg
}
