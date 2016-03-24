package irc

import (
	"testing"
)

func TestNilConfigConstructor(t *testing.T) {
	cfg := NewConfig(nil)
	if cfg == nil {
		t.Error("Got nil config, non-nil expected")
	}
	if cfg.User != nil {
		t.Error("Got non-nil config user with nil parameter, nil expected")
	}
}

func TestConfigConstructor(t *testing.T) {
	user := &User{
		Username: "test",
		Realname: "test",
	}
	cfg := NewConfig(user)
	if cfg.User == nil {
		t.Error("Config user mismatch, got nil, expected non-nil")
	}
	if cfg.User.Username != "test" {
		t.Errorf("Config user username mismatch, got '%s', expected '%s'", cfg.User.Username, "test")
	}
	if cfg.User.Realname != "test" {
		t.Errorf("Config user realname mismatch, got '%s', expected '%s'", cfg.User.Realname, "test")
	}
}
