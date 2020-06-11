package svm

type HostCtx map[uint32][]byte

func NewHostCtx() HostCtx {
	return make(HostCtx)
}

func (h HostCtx) Encode() []byte {
	if len(h) == 0 {
		var size int
		size += 4 // proto version
		size += 2 // #fields
		return make([]byte, size, size)
	}

	panic("not implemented")
}
