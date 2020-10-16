package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/coreos/go-systemd/activation"
)

// Listen creates a new listener.
func Listen(address string) (net.Listener, error) {
	if !strings.HasPrefix(s.address, "systemd:") {
		return net.Listen("tcp", address)
	}

	listeners, err := activation.ListenersWithNames()
	if err != nil {
		return nil, err
	}

	name := s.address[8:]
	listener, ok := listeners[name]

	if !ok {
		return nil, fmt.Errorf("listen systemd %s: socket not found", name)
	}

	return listener[0], nil
}
