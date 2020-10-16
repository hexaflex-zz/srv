package main

import "net"

// Listen creates a new listener.
func Listen(address string) (net.Listener, error) {
	return net.Listen("tcp", address)
}
