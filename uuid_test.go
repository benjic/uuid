package uuid

import (
	"bytes"
	"testing"
)

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
