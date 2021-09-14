package svm

type ValidateErrorKind byte
type RuntimeErrorKind byte

const (
	ParseError    ValidateErrorKind = 0
	ProgramError  ValidateErrorKind = 1
	FixedGasError ValidateErrorKind = 2
)

const (
	OOG                  RuntimeErrorKind = 0
	TemplateNotFound     RuntimeErrorKind = 1
	AccountNotFound      RuntimeErrorKind = 2
	CompilationFailed    RuntimeErrorKind = 3
	InstantiationFailed  RuntimeErrorKind = 4
	FuncNotFound         RuntimeErrorKind = 5
	FuncFailed           RuntimeErrorKind = 6
	FuncNotAllowed       RuntimeErrorKind = 7
	FuncInvalidSignature RuntimeErrorKind = 8
)

type ValidateError struct {
	Kind    ValidateErrorKind
	Message string
}

type RuntimeError struct {
	Kind     RuntimeErrorKind
	Target   Address
	Function string
	Template TemplateAddr
	Message  string
}
