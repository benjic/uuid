package uuid

import (
	"bytes"
	"testing"
)

func resetDefaultGenerate() {
	defaultGenerator = nil
}

func TestDefaultGenerator(t *testing.T) {
	resetDefaultGenerate()

	uuid := Generate()

	if uuid.Version() != 4 {
		t.Errorf("Expected a version 4 uuid from the default generator; got %d", uuid.Version())
	}
}

func TestDefaultGeneratorPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected Generate to panic if defaultConfiguration is invalid")
		}
	}()

	resetDefaultGenerate()
	defaultConfiguration = Configuration{Version: 100}

	// Attempt to generate with an invalided defualt configuration
	Generate()
}

func TestUUIDString(t *testing.T) {
	cases := []struct {
		uuid     UUID
		expected string
	}{
		{
			[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			"00000000-0000-0000-0000-000000000000",
		},
		{
			[]byte{248, 29, 79, 174, 125, 236, 17, 208, 167, 101, 0, 160, 201, 30, 107, 246},
			"f81d4fae-7dec-11d0-a765-00a0c91e6bf6",
		},
	}

	for _, c := range cases {
		got := c.uuid.String()

		if got != c.expected {
			t.Errorf("Expected %+v bytes to return %s; got %s", []byte(c.uuid), c.expected, got)
		}
	}
}

func TestUUIDParse(t *testing.T) {
	cases := []struct {
		err      error
		expected UUID
		uuid     string
	}{
		{
			nil,
			[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			"00000000-0000-0000-0000-000000000000",
		},
		{
			nil,
			[]byte{248, 29, 79, 174, 125, 236, 17, 208, 167, 101, 0, 160, 201, 30, 107, 246},
			"f81d4fae-7dec-11d0-a765-00a0c91e6bf6",
		},
	}

	for _, c := range cases {
		got, err := Parse(c.uuid)

		if err != c.err {
			t.Errorf("Expected %s error when parsing; got %s", c.err, err)
		}

		if !bytes.Equal(got, c.expected) {
			t.Errorf("Expected %s string to return:\n want:%+v\n  got:%+v", c.uuid, []byte(c.expected), []byte(got))
		}
	}
}

func TestUUIDVersion(t *testing.T) {
	cases := []struct {
		uuid    string
		version Version
	}{
		{"342ce962-286b-11e6-b67b-9e71128cae77", 1},
		{"a1322ba2-939d-2c50-91f9-7f1b6465941f", 2},
		{"d0021c5b-595b-362b-ac92-0ac7f3faf74f", 3},
		{"75c02ed9-f82f-4c68-91a9-9b0b6c2dd698", 4},
		{"031ea7d5-c569-5d37-b6c2-9d1a58a96316", 5},
	}

	for _, c := range cases {
		uuid, err := Parse(c.uuid)

		if err != nil {
			t.Errorf("Did not expect error %s when parsing test case UUID", err)
		}

		got := uuid.Version()

		if got != c.version {
			t.Errorf("Expected %s to have version %d; Got %d", c.uuid, c.version, got)
		}
	}
}
