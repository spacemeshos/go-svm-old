package svm

const SvmVersion = 0

type Section interface {
	encode(bs []byte) ([]byte, error)
	length() int
	kind() SectionKind
}

type SectionKind uint16

const (
	CodeSection  SectionKind = 1
	DataSection  SectionKind = 2
	CtorsSection SectionKind = 3
)

type CodeKind uint16

const (
	CodeKindWasm CodeKind = 1
)

type CodeFlags uint64

const (
	CodeExecFlags CodeFlags = 0x01
)

type CodeGasMode uint16

const (
	GasModeFixed    CodeGasMode = 1
	GasModeMetering CodeGasMode = 2
)

type WasmFixedGasCode []byte

func (cs WasmFixedGasCode) kind() SectionKind {
	return CodeSection
}

func (cs WasmFixedGasCode) length() int {
	return 2 + 8 + 8 + 4 + 4 + len(cs)
}

func (cs WasmFixedGasCode) encode(bs []byte) ([]byte, error) {
	bs = Encode16be(bs, uint16(CodeKindWasm))
	bs = Encode64be(bs, uint64(CodeExecFlags))
	bs = Encode64be(bs, uint64(GasModeFixed))
	bs = Encode32be(bs, SvmVersion)
	bs = Encode32be(bs, uint32(len(cs)))
	bs = append(bs, cs...)
	return bs, nil
}

type Ctors []string

func (cs Ctors) kind() SectionKind {
	return CtorsSection
}

func (cs Ctors) length() int {
	l := 1
	for _, v := range cs {
		l += 1 + len(v)
	}
	return l
}

func (cs Ctors) encode(bs []byte) ([]byte, error) {
	bs = append(bs, byte(len(cs)))
	for _, v := range cs {
		bs = append(bs, byte(len(v)))
		bs = append(bs, v...)
	}
	return bs, nil
}

type LayoutKind uint16

const (
	Fixed LayoutKind = 1
)

type FixedLayout []uint32

func (fl FixedLayout) kind() SectionKind {
	return DataSection
}

func (fl FixedLayout) encode(bs []byte) ([]byte, error) {
	bs = Encode16be(bs, uint16(1))
	bs = Encode16be(bs, uint16(Fixed))
	bs = Encode16be(bs, uint16(len(fl)))
	bs = Encode32be(bs, 0)
	for _, v := range fl {
		bs = Encode16be(bs, uint16(v))
	}
	return bs, nil
}

func (fl FixedLayout) length() int {
	return 2
}

const previewSectionLength = 8 // ?
func encodePreview(bs []byte, s Section) ([]byte, error) {
	bs = Encode16be(bs, uint16(s.kind()))
	bs = Encode32be(bs, uint32(s.length()))
	return bs, nil
}
