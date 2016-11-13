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

// NewGenerator creates a new UUID generator which operates with the given
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
	uuid := make([]byte, 16)
	g.reader.Read(uuid)
	applyFlags(uuid, g.Version)
	return uuid
}

func applyFlags(uuid UUID, version Version) {
	// Apply flags
	uuid[6] = byte(version<<4) | (0x0f & uuid[6])
	uuid[8] = 0xBF & (0x80 | uuid[8])
}
