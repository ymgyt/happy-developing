package errors

import "bytes"

// Code -
type Code int

const (
	// OK -
	OK Code = iota
	// InvalidInputErr -
	InvalidInputErr
	// UndefinedCode -
	UndefinedCode
)

func (c Code) String() (s string) {
	switch c {
	case OK:
		s = "ok"
	case InvalidInputErr:
		s = "invalidInput"
	default:
		s = "undefined"
	}
	return s
}

// GetCode -
func GetCode(err error) Code {
	if err == nil {
		return OK
	}
	appErr, ok := err.(*Error)
	if !ok {
		return UndefinedCode
	}
	return appErr.Code
}

// New -
func New(code Code, msg string) *Error {
	return &Error{Code: code, Message: msg}
}

// Wrap -
func Wrap(code Code, msg string, cause error) *Error {
	e := New(code, msg)
	e.Cause = cause
	return e
}

// InvalidInput -
func InvalidInput(msg string, cause error) *Error {
	return createOrWrap(InvalidInputErr, msg, cause)
}

func createOrWrap(code Code, msg string, cause error) *Error {
	if cause == nil {
		return New(code, msg)
	}
	return Wrap(code, msg, cause)
}

// Error -
type Error struct {
	Code    Code
	Message string
	Cause   error
}

// Error -
func (e *Error) Error() string {
	var b bytes.Buffer
	if e.Code != OK {
		b.WriteString(e.Code.String() + ":")
	}
	b.WriteString(e.Message)

	if e.Cause != nil {
		b.WriteString("/")
		b.WriteString(e.Cause.Error())
	}
	return b.String()
}
