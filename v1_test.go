package uuid

import (
	"bytes"
	"errors"
	"net"
	"testing"
	"time"
)

func TestTimeFactoryResolution(t *testing.T) {
	// Ensure the timeFactory will receive the timestamp
	testTime := func() time.Time { return time.Unix(0, 0) }

	tFactory := createTimeFactory(testTime)

	cases := []struct {
		timestamp []byte
	}{
		{[]byte{1, 178, 29, 210, 19, 129, 64, 0}},
		{[]byte{1, 178, 29, 210, 19, 129, 64, 1}},
		{[]byte{1, 178, 29, 210, 19, 129, 64, 2}},
	}

	for _, c := range cases {
		bs := tFactory()

		if !bytes.Equal(bs, c.timestamp) {
			t.Error("Expected timeFactory to return incremented result if clock resolves same time")
		}
	}
}

func TestV1ReaderRead(t *testing.T) {
	// Create a set of unique bytes returned by factories
	timeFactory := func() []byte { return []byte{0, 1, 2, 3, 4, 5, 6, 7} }
	node := []byte{10, 11, 12, 13, 14, 15}
	clockSequence := []byte{255, 255}

	testV1Reader := &v1Reader{clockSequence, node, timeFactory}

	output := make([]byte, 16)

	count, err := testV1Reader.Read(output)

	if err != nil {
		t.Error("Did not expect error %s", err)
	}

	if count != 16 {
		t.Error("Read did not encounter the correct number of bytes; Wanted 16 got %d", count)
	}

	cases := []struct {
		label    string
		start    int
		end      int
		expected []byte
	}{
		{"timestep_low", 0, 4, []byte{4, 5, 6, 7}},
		{"timestep_mid", 4, 6, []byte{2, 3}},
		{"timestep_high", 6, 8, []byte{0, 1}},
		{"clock_sequence", 8, 10, []byte{255, 255}},
		{"node", 10, 16, []byte{10, 11, 12, 13, 14, 15}},
	}

	for _, c := range cases {
		got := output[c.start:c.end]
		if !bytes.Equal(got, c.expected) {
			t.Errorf("Expected %s bytes to be %+v, got %+v", c.label, c.expected, got)
		}
	}
}

func TestNewV1Reader(t *testing.T) {
	firstErr := errors.New("Failed random bytes invocation")
	secondErr := errors.New("Failed to generate bytes for clock sequence")

	failOn := func(invocation int) func([]byte) (int, error) {
		count := 0
		return func(bs []byte) (int, error) {
			count++
			switch count {
			case 1:
				if invocation == 1 {
					return 0, firstErr
				}

				copy(bs, []byte{1, 1, 1, 1, 1, 1})
				return 6, nil

			case 2:
				if invocation == 2 {
					return 0, secondErr
				}

				copy(bs, []byte{2, 2})
				return 2, nil
			default:
				panic("The failOn function should only be used once")
			}
		}
	}

	cases := []struct {
		err                   error
		expectedClockSequence []byte
		expectedNode          []byte
		ifts                  []net.Interface
		randomRead            func([]byte) (int, error)
	}{
		{firstErr, []byte{0, 0}, []byte{0, 0, 0, 0, 0, 0}, []net.Interface{}, failOn(1)},
		{secondErr, []byte{0, 0}, []byte{1, 1, 1, 1, 1, 1}, []net.Interface{}, failOn(2)},
		{nil, []byte{2, 2}, []byte{1, 1, 1, 1, 1, 1}, []net.Interface{}, failOn(3)},
	}

	for _, c := range cases {
		got, err := newV1Reader(c.ifts, c.randomRead)

		if err != c.err {
			t.Errorf("Expected error %s; got %s", c.err, err)
		}

		if got != nil {
			if !bytes.Equal(got.clockSequence, c.expectedClockSequence) {
				t.Errorf("Expected clockSequence %+v; got %+v", c.expectedClockSequence, got.clockSequence)
			}

			if !bytes.Equal(got.node, c.expectedNode) {
				t.Errorf("Expected node %+v; got %+v", c.expectedNode, got.node)
			}
		}
	}
}

func TestNodeValue(t *testing.T) {
	defaultNode := func(bs []byte) (int, error) {
		copy(bs, []byte{1, 1, 1, 1, 1, 1})
		return 2, nil
	}

	testErr := errors.New("This is a test error")
	errDefaultNode := func(bs []byte) (int, error) {
		return 0, testErr
	}

	cases := []struct {
		defaultNode func([]byte) (int, error)
		err         error
		expected    []byte
		ifts        []net.Interface
		message     string
	}{
		{defaultNode, nil, []byte{1, 1, 1, 1, 1, 1}, []net.Interface{}, "Expected default value %+v when interfaces is empty; got %+v"},
		{defaultNode, nil, []byte{1, 1, 1, 1, 1, 1}, nil, "Expected default value %+v when interfaces is nil; got %+v"},
		{
			defaultNode,
			nil,
			[]byte{1, 1, 1, 1, 1, 1},
			[]net.Interface{net.Interface{}},
			"Expected default value %+v when interfaces have no HardwareAddr; got %+v",
		},
		{
			defaultNode,
			nil,
			[]byte{2, 2, 2, 2, 2, 2},
			[]net.Interface{
				net.Interface{HardwareAddr: []byte{2, 2, 2, 2, 2, 2}},
			},
			"Expected hardware interface value %+v, got %+v",
		},
		{
			defaultNode,
			nil,
			[]byte{2, 2, 2, 2, 2, 2},
			[]net.Interface{
				net.Interface{HardwareAddr: []byte{2, 2, 2, 2, 2, 2}},
				net.Interface{HardwareAddr: []byte{3, 3, 3, 3, 3, 3}},
			},
			"Expected first hardware interface value %+v, got %+v",
		},
		{
			errDefaultNode,
			testErr,
			[]byte{0, 0, 0, 0, 0, 0},
			nil,
			"Expected returned error %s; got %s",
		},
	}

	for _, c := range cases {
		got, err := findOrDefaultNode(c.ifts, c.defaultNode)

		if err != c.err {
			t.Errorf("Expected error %s; got %s", c.err, err)
			break
		}

		if !bytes.Equal(got, c.expected) {
			t.Errorf(c.message, c.expected, got)
		}
	}
}
