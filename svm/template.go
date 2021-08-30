package svm

// SvmVersion is a current SVM version
const SvmVersion = 0

type section interface {
	encode(bs []byte) ([]byte, error)
	length() int
	kind() sectionKind
}

type sectionKind uint16

const (
	// CodeSection is a kind of section for a code definition
	CodeSection  sectionKind = 1
	// DataSection is a kind of section for data layouts
	DataSection  sectionKind = 2
	// CtorsSection is a kind of section for constructors
	CtorsSection sectionKind = 3
)

type codeKind uint16

const (
	// WasmCode defines code as wasm assembly
	WasmCode codeKind = 1
)

type codeFlags uint64

const (
	// CodeExecFlags is the default code flags
	CodeExecFlags codeFlags = 0x01
)

type codeGasMode uint16

const (
	// GasModeFixed is a default mode for gas calculation for the code
	GasModeFixed    codeGasMode = 1
)

// WasmFixedGasCode is the abstraction for wasm code assembly with fixed gas calculation
type WasmFixedGasCode []byte

func (cs WasmFixedGasCode) kind() sectionKind {
	return CodeSection
}

func (cs WasmFixedGasCode) length() int {
	return 2 + 8 + 8 + 4 + 4 + len(cs)
}

func (cs WasmFixedGasCode) encode(bs []byte) ([]byte, error) {
	bs = encode16be(bs, uint16(WasmCode))
	bs = encode64be(bs, uint64(CodeExecFlags))
	bs = encode64be(bs, uint64(GasModeFixed))
	bs = encode32be(bs, SvmVersion)
	bs = encode32be(bs, uint32(len(cs)))
	bs = append(bs, cs...)
	return bs, nil
}

// Ctors is a constructors section abstraction
type Ctors []string

func (cs Ctors) kind() sectionKind {
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
	// Fixed is a fixed data layout
	Fixed LayoutKind = 1
)

// FixedLayot is a data section abstraction
type FixedLayout []uint32

func (fl FixedLayout) kind() sectionKind {
	return DataSection
}

func (fl FixedLayout) encode(bs []byte) ([]byte, error) {
	bs = encode16be(bs, uint16(1))
	bs = encode16be(bs, uint16(Fixed))
	bs = encode16be(bs, uint16(len(fl)))
	bs = encode32be(bs, 0)
	for _, v := range fl {
		bs = encode16be(bs, uint16(v))
	}
	return bs, nil
}

func (fl FixedLayout) length() int {
	return 2
}

const previewSectionLength = 8 // ?
func encodePreview(bs []byte, s section) ([]byte, error) {
	bs = encode16be(bs, uint16(s.kind()))
	bs = encode32be(bs, uint32(s.length()))
	return bs, nil
}
