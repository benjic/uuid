// Package uuid is a simple library for generating RFC4122 compliant universally unique
// identifiers.
package uuid

import (
	"crypto/rand"
	"fmt"
)

const (
	stringFormat = "%8x-%4x-%4x-%4x-%12x"
)

var defaultConfiguration Configuration = Configuration{
	Version:      4,
	RandomReader: rand.Read,
}
var defaultGenerator *Generator

// Generate provides a Version 4 random UUID
func Generate() UUID {
	var err error

	if defaultGenerator == nil {
		defaultGenerator, err = NewGenerator(defaultConfiguration)

		if err != nil {
			panic(fmt.Sprintf("Unable to create default generator: %s", err))
		}
	}

	return defaultGenerator.Generate()
}

// A UUID is a unique identifier.
type UUID []byte

// String provides a string representation of a UUID using the established
// dashed format.
func (id UUID) String() string {
	bs := []byte(id)
	return fmt.Sprintf(stringFormat, bs[0:4], bs[4:6], bs[6:8], bs[8:10], bs[10:])
}

// Parse returns a UUID from a given string
func Parse(id string) (uuid UUID, err error) {
	var a, b, c, d, e []byte

	_, err = fmt.Sscanf(id, stringFormat, &a, &b, &c, &d, &e)

	uuid = append(append(append(append(a, b...), c...), d...), e...)

	return uuid, err
}

// Version of the UUID
func (id UUID) Version() Version { return Version((id[6] & 0xF0) >> 4) }
