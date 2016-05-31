package uuid

import "testing"

type testWriter struct{ wasWritten bool }

func (w *testWriter) Write(buf []byte) (int, error) {
	w.wasWritten = true

	return 16, nil
}

type byteWriter struct {
	index int
	done  bool
	count uint8
}

func (w *byteWriter) Write(buf []byte) (int, error) {

	if w.count == 0xff {
		w.done = true
		return 0, nil
	}

	buf[w.index] = byte(w.count)
	w.count += 1

	return 16, nil
}

func testGeneratorFactory(version Version) *Generator {
	return &Generator{&testWriter{}, version}
}

func byteWriterGeneratorFactory(index int, version Version) (*Generator, *byteWriter) {
	w := &byteWriter{index, false, 0}
	gen := &Generator{w, version}

	return gen, w
}

func TestNewGeneratorInvalidVersion(t *testing.T) {
	g, err := NewGenerator(-1)

	if err != errUnknownVersion {
		t.Errorf("Expected error %s with version %d; got %s", errUnknownVersion, -1, err)
	}

	if g != nil {
		t.Errorf("Expected invalid version to return a nil generator")
	}
}

func TestGeneratorVariant(t *testing.T) {
	gen, w := byteWriterGeneratorFactory(8, 0)

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

		gen, w := byteWriterGeneratorFactory(6, version)

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
