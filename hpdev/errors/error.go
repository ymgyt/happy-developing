package errors

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Code -
type Code int

const (
	// OK -
	OK Code = iota
	// InvalidInputErr -
	InvalidInputErr
	// AlreadyExistsErr -
	AlreadyExistsErr
	// UndefinedCode -
	UndefinedCode
)

func (c Code) String() (s string) {
	switch c {
	case OK:
		s = "ok"
	case InvalidInputErr:
		s = "invalidInput"
	case AlreadyExistsErr:
		s = "alreadyExists"
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

// StatusCode -
func StatusCode(err error) (code int) {
	if err == nil {
		return http.StatusOK
	}
	appErr, ok := err.(*Error)
	if !ok {
		return http.StatusInternalServerError
	}
	switch appErr.Code {
	case OK:
		code = http.StatusOK
	case InvalidInputErr, AlreadyExistsErr:
		code = http.StatusBadRequest
	default:
		code = http.StatusInternalServerError
	}
	return code
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

// AlreadyExists -
func AlreadyExists(msg string, cause error) *Error {
	return createOrWrap(AlreadyExistsErr, msg, cause)
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

// MarshalJSON -
func (e *Error) MarshalJSON() ([]byte, error) {
	m := map[string]string{
		"code":    e.Code.String(),
		"message": e.Message,
	}
	if e.Cause != nil {
		m["message"] = e.Cause.Error()
	}

	return json.Marshal(&m)
}
