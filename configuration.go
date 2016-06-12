package uuid

import "net"

type Version uint8

// A Configuration provides sources of things.
type Configuration struct {
	Version      Version
	Interfaces   []net.Interface
	RandomReader func([]byte) (int, error)
}
