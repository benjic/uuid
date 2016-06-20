package uuid

import "net"

// A Version dictates the generation of an UUID
type Version uint8

// A Configuration provides sources of things.
type Configuration struct {
	Version      Version
	Interfaces   []net.Interface
	RandomReader func([]byte) (int, error)
}
