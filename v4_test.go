package uuid

import (
	"net"
	"testing"
)

type readProxy struct {
	didCall bool
	payload []byte
	err     error
}

func (rp *readProxy) Read(bs []byte) (int, error) {
	rp.didCall = true
	copy(bs, rp.payload)

	return len(bs), rp.err
}

func TestVersion4Reader(t *testing.T) {
	rp := &readProxy{}
	config := Configuration{
		4,
		[]net.Interface{},
		rp.Read,
	}

	gen, err := NewGenerator(config)

	if err != nil {
		t.Errorf("Did not expect an error; Got %s", err)
	}

	gen.Generate()

	if !rp.didCall {
		t.Errorf("Expected generator to call underlying random proxy")
	}

}
