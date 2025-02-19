package errs

import "fmt"

type Extras map[string]any

type Error struct {
	StatusCode int     `json:"-" swaggerignore:"true"`
	ErrorCode  string  `json:"code" example:"INVALID_REQUEST"`
	Message    string  `json:"message" example:"invalid request: some or all request parameters are invalid"`
	Extras     *Extras `json:"-"`
}

func New(statusCode int, errorCode, message string) *Error {
	return &Error{
		StatusCode: statusCode,
		ErrorCode:  errorCode,
		Message:    message,
	}
}

func NewImmutable(statusCode int, errorCode, message string) Error {
	return Error{
		StatusCode: statusCode,
		ErrorCode:  errorCode,
		Message:    message,
	}
}

func (e Error) Msg(format string, parts ...any) *Error {
	e.Message = fmt.Sprintf(format, parts...)
	return &e
}

func (e Error) WithExtras(extras Extras) *Error {
	e.Extras = &extras
	return &e
}

func NewInvalidViolations(violations any) *Error {
	// copy ErrInvalidRequest as e
	e := *ErrInvalidReq
	e.Extras = &Extras{
		"violations": violations,
	}
	return &e
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode, e.Message)
}
