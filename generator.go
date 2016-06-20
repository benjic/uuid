package uuid

import (
	"errors"
	"io"
)

var (
	errUnknownVersion = errors.New("Unknown version")
)

// A Generator is a factory for creating new UUIDs
type Generator struct {
	reader  io.Reader
	Version Version
}

// NewGenerator creates a new  UUID generator which operates with the given
// configuration.
func NewGenerator(configuration Configuration) (g *Generator, err error) {

	switch configuration.Version {
	case 1:
		reader, err := newV1Reader(configuration.Interfaces, configuration.RandomReader)

		if err != nil {
			return nil, err
		}

		return &Generator{reader, configuration.Version}, nil
	case 4:
		return &Generator{v4Reader{configuration.RandomReader}, configuration.Version}, nil
	}

	return nil, errUnknownVersion
}

// Generate produces a new UUID that reflects the configuration of the
// generator.
func (g *Generator) Generate() UUID {
	bs := make([]byte, 16)

	g.reader.Read(bs)
	// TODO: communicate error

	// Apply flags
	bs[6] = byte(g.Version<<4) | (0x0f & bs[6])
	bs[8] = 0xBF & (0x80 | bs[8])

	return bs
}
