package uuid

import (
	"errors"
	"io"
)

var (
	errUnknownVersion = errors.New("Unknown version")
)

// A Version identifies how the generator generates UUIDs
type Version int8

// A Generator is a factory for creating new UUIDs
type Generator struct {
	writer  io.Writer
	Version Version
}

// NewGenerator creates a new  UUID generator configured for the given version.
func NewGenerator(version Version) (g *Generator, err error) {

	// TODO: Implement version based writers
	switch version {
	}

	return nil, errUnknownVersion
}

// Generate produces a new UUID that reflects the configuration of the
// generator.
func (g *Generator) Generate() UUID {
	bs := make([]byte, 16)

	g.writer.Write(bs)

	// Apply flags
	bs[6] = byte(g.Version<<4) | (0x0f & bs[6])
	bs[8] = 0xBF & (0x80 | bs[8])

	return bs
}
