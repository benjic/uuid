package uuid

import "testing"
import "bytes"

var aNamespaceUUID UUID = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
var bNamespaceUUID UUID = []byte{16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}

func TestNewNamespaceGenerator(t *testing.T) {

	unsupportedVersions := []Version{1, 2, 4}
	supportedVersions := []Version{3, 5}

	for _, version := range unsupportedVersions {
		config := Configuration{Version: version}
		g, err := NewNamespaceGenerator(aNamespaceUUID, config)

		if err != errUnknownVersion {
			t.Error("Expected NewNamespaceGenerator to return error with non namespaced versions")
		}

		if g != nil {
			t.Error("Expected nil return value when error is returned")
		}
	}

	for _, version := range supportedVersions {
		config := Configuration{Version: version}
		g, err := NewNamespaceGenerator(aNamespaceUUID, config)

		if err != nil {
			t.Error("Expected supported version to not return an error")
		}

		if g.hashFactory == nil {
			t.Error("Did not expect null hashing factory")
		}

		if !bytes.Equal(aNamespaceUUID, g.Namespace) {
			t.Error("Expected given namespace to be attributed to the generator")
		}
	}

}

func TestNamespacedGenerate(t *testing.T) {
	testName := []byte("I love UUIDs")
	secondName := []byte("I hate UUIDs")
	versions := []Version{3, 5}

	for _, version := range versions {
		aGenerator, aErr := NewNamespaceGenerator(aNamespaceUUID, Configuration{Version: version})

		if aErr != nil {
			t.Error("Did not expect error when creating a new V3 Generator")
		}

		bGenerator, bErr := NewNamespaceGenerator(bNamespaceUUID, Configuration{Version: version})

		if bErr != nil {
			t.Error("Did not expect error when creating new Namespaced generator")
		}

		// Same name in a given namespace should return the same UUID
		first := aGenerator.Generate(testName)
		second := aGenerator.Generate(testName)

		if !bytes.Equal(first, second) {
			t.Error("Expected the same name in a single namespace to return same UUID")
		}

		// Differing names should produce differing UUIDs
		third := aGenerator.Generate(secondName)
		if bytes.Equal(first, third) {
			t.Error("Expected differing names in single namespace to return differing UUIDs")
		}

		fourth := bGenerator.Generate(testName)
		if bytes.Equal(first, fourth) {
			t.Error("Expected same name in diferrent namespaces to return differing UUIDs")
		}
	}
}
