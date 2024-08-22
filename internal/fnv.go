package fnv

type fnv64a uint64

func New64a() fnv64a {
	fnv := fnv64a(0xcbf29ce484222325)
	return fnv
}

func (h *fnv64a) Write(p []byte) (n int, err error) {
	for _, b := range p {
		h.WriteByte(b)
	}
	return len(p), nil
}

func (h *fnv64a) WriteByte(b byte) error {
	*h ^= fnv64a(b)
	*h *= 0x100000001b3
	return nil
}

func (h *fnv64a) WriteUint64(v uint64) {
	h.WriteByte(byte(0xff & v))
	h.WriteByte(byte(0xff & (v >> 8)))
	h.WriteByte(byte(0xff & (v >> 16)))
	h.WriteByte(byte(0xff & (v >> 24)))
	h.WriteByte(byte(0xff & (v >> 32)))
	h.WriteByte(byte(0xff & (v >> 40)))
	h.WriteByte(byte(0xff & (v >> 48)))
	h.WriteByte(byte(0xff & (v >> 56)))
}

func (h fnv64a) Sum() uint64 {
	return uint64(h)
}
