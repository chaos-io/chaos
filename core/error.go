package core

import (
	"errors"
	"fmt"
)

type Error struct {
	Code    *ErrorCode `json:"code"`
	Message string     `json:"message"`
	Details []any      `json:"details"`
}

func NewError(code *ErrorCode, message string, arguments ...any) *Error {
	msg := message
	if len(arguments) > 0 {
		msg = fmt.Sprintf(message, arguments...)
	}

	return &Error{
		Code:    code,
		Message: msg,
	}
}

func NewErrorFrom(code int32, message string, arguments ...any) *Error {
	err := &Error{}
	if ec, ok := errorCodeIndex[code]; ok {
		err.Code = ec
	} else {
		err.Code = &ErrorCode{Code: code}
	}

	err.Message = message
	if len(arguments) > 0 {
		err.Message = fmt.Sprintf(message, arguments...)
	}

	return err
}

func (e *Error) Is(err error) bool {
	return IsError(err)
}

func IsError(err error) bool {
	return errors.Is(err, &Error{})
}

func AsError(err error) *Error {
	e := &Error{}
	if errors.As(err, &e) {
		return e
	}
	return nil
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if len(e.Message) == 0 {
		return e.Code.Name
	}
	return e.Message
}

func (e *Error) StatusCode() int {
	if e == nil {
		return 200
	}
	if e.Code.HttpStatusCode > 0 {
		return int(e.Code.HttpStatusCode)
	}
	return int(e.Code.Code)
}

func (e *Error) AddDetail(detail any) *Error {
	if e != nil {
		e.Details = append(e.Details, detail)
	}
	return e
}
