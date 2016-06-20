package uuid

import (
	"errors"
	"net"
	"testing"
)

type testReader struct{ wasWritten bool }

func (w *testReader) Read(buf []byte) (int, error) {
	w.wasWritten = true

	return 16, nil
}

type byteReader struct {
	index int
	done  bool
	count uint8
}

func (w *byteReader) Read(buf []byte) (int, error) {

	if w.count == 0xff {
		w.done = true
		return 0, nil
	}

	buf[w.index] = byte(w.count)
	w.count += 1

	return 16, nil
}

func testGeneratorFactory(version Version) *Generator {
	return &Generator{&testReader{}, version}
}

func byteReaderGeneratorFactory(index int, version Version) (*Generator, *byteReader) {
	w := &byteReader{index, false, 0}
	gen := &Generator{w, version}

	return gen, w
}

func newTestConfiguration(version Version) Configuration {
	return Configuration{
		version,
		[]net.Interface{},
		func(bs []byte) (int, error) { return 0, nil },
	}
}

func TestNewGeneratorInvalidVersion(t *testing.T) {
	config := newTestConfiguration(255)
	g, err := NewGenerator(config)

	if err != errUnknownVersion {
		t.Errorf("Expected error %s with version %d; got %s", errUnknownVersion, -1, err)
	}

	if g != nil {
		t.Errorf("Expected invalid version to return a nil generator")
	}
}

func TestNewGeneratorFailedVersion1(t *testing.T) {
	expectedErr := errors.New("Failed to read random bytes")

	config := Configuration{
		1,
		[]net.Interface{},
		func(bs []byte) (int, error) { return 0, expectedErr },
	}

	g, err := NewGenerator(config)

	if err != expectedErr {
		t.Errorf("Expected faulted random source to return an error %s; got %s", expectedErr, err)
	}

	if g != nil {
		t.Errorf("Expected a nil generator when an error is returned")
	}
}

func TestNewGeneratorValidVersions(t *testing.T) {
	versions := []Version{1, 4}

	for _, version := range versions {
		config := newTestConfiguration(version)
		_, err := NewGenerator(config)

		if err != nil {
			t.Errorf("Did not expect generator to return error for valid version; got %s", err)
		}
	}
}

func TestGeneratorVariant(t *testing.T) {
	gen, w := byteReaderGeneratorFactory(8, 0)

	for !w.done {
		uuid := gen.Generate()

		varaint := uuid[8] & 0xC0

		if varaint != byte(0x80) {
			t.Errorf("Expected any value of varaint field to begin with 0x04; Got %x", varaint)
		}

	}
}

func TestGeneratorVersion(t *testing.T) {
	versions := []Version{1, 2, 3, 4, 5}

	for _, version := range versions {

		gen, w := byteReaderGeneratorFactory(6, version)

		for !w.done {
			uuid := gen.Generate()

			got := uuid.Version()

			if got != version {
				t.Errorf("Expected any value of version field to equal generator version %d; Got %d", version, got)
			}
		}
	}
}

func BenchmarkGeneratorGenerate(b *testing.B) {
	gen := testGeneratorFactory(1)

	for i := 0; i < b.N; i++ {
		gen.Generate()
	}
}
