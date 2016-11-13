package uuid

import (
	"crypto/md5"
	"crypto/sha1"
	"hash"
)

// A NamespaceGenerator allows consumers to generate UUIDs for a given namespace
// and name.
type NamespaceGenerator struct {
	hashFactory func() hash.Hash
	Namespace   UUID
	Version     Version
}

// NewNamespaceGenerator creates a new UUID generator for the given namespace
// and configuration.
func NewNamespaceGenerator(namespace UUID, configuration Configuration) (g *NamespaceGenerator, err error) {
	switch configuration.Version {
	case 3:
		return &NamespaceGenerator{md5.New, namespace, configuration.Version}, nil
	case 5:
		return &NamespaceGenerator{sha1.New, namespace, configuration.Version}, nil
	}
	return nil, errUnknownVersion
}

// Generate procues a new UUID that refelcts the given
func (g *NamespaceGenerator) Generate(name []byte) UUID {
	h := g.hashFactory()
	data := append(g.Namespace, name...)
	uuid := h.Sum(data)
	applyFlags(uuid, g.Version)
	return uuid
}
