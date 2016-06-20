package uuid

type v4Reader struct {
	RandomReader func([]byte) (int, error)
}

func (w v4Reader) Read(bs []byte) (int, error) { return w.RandomReader(bs) }
