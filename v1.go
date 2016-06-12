package uuid

import (
	"bytes"
	"encoding/binary"
	"net"
	"time"
)

const (
	// Good thing calendars make sense
	incrementsInASecond  uint64 = 1e7
	incrementsInDay      uint64 = 60 * 60 * 24 * incrementsInASecond
	daysUntilEpoch       uint64 = 77 + (365 * 387) + (387 / 4) - 1
	incrementsUntilEpoch uint64 = daysUntilEpoch * incrementsInDay
)

type interfaces []net.Interface

type v1Reader struct {
	clockSequence []byte
	node          []byte
	timeFactory   func() []byte
}

func createTimeFactory(currentTime func() time.Time) func() []byte {
	var last time.Time
	var count uint64

	return func() []byte {
		buf := new(bytes.Buffer)

		now := currentTime().UTC()
		incrementsSinceEpoch := uint64(now.UnixNano() / 100)

		if now.Equal(last) {
			count++
		}

		// TODO: Make this thread safe
		last = now

		binary.Write(buf, binary.BigEndian, incrementsUntilEpoch+incrementsSinceEpoch+count)
		return buf.Bytes()
	}
}

func findOrDefaultNode(ifts interfaces, defaultNode func([]byte) (int, error)) (net.HardwareAddr, error) {

	// Find any interface with a hardware address in the given list
	for _, ift := range ifts {
		if ift.HardwareAddr != nil {
			return ift.HardwareAddr, nil
		}
	}

	// Use a default value if no interfaces are found
	bs := make([]byte, 6)
	_, err := defaultNode(bs)

	return bs, err
}

func newV1Reader(ifts interfaces, randomRead func([]byte) (int, error)) (*v1Reader, error) {

	hardwareAddr, err := findOrDefaultNode(ifts, randomRead)

	if err != nil {
		return nil, err
	}

	timeFactory := createTimeFactory(time.Now)

	clockSequence := make([]byte, 2)
	_, err = randomRead(clockSequence)

	if err != nil {
		return nil, err
	}

	return &v1Reader{
		clockSequence,
		hardwareAddr,
		timeFactory,
	}, err
}

func (w *v1Reader) Read(bs []byte) (int, error) {
	time := w.timeFactory()

	// Rearrange the time bits
	copy(bs[0:4], time[4:8])
	copy(bs[4:6], time[2:4])
	copy(bs[6:8], time[0:2])

	copy(bs[8:10], w.clockSequence)
	copy(bs[10:16], w.node)

	return 16, nil
}
