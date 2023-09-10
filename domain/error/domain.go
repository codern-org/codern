package errs

import (
	"errors"
	"fmt"
)

type DomainError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`

	// Underlying error, if any
	Err error `json:"-"`
}

func New(code int, message string, args ...interface{}) *DomainError {
	var err error
	var ok bool
	i := 0

	if len(args) > 0 {
		// Assume the last args is a wrapped error
		err, ok = args[len(args)-1].(error)
		if ok {
			i = 1
		}
	}
	message = fmt.Sprintf(message, args[:len(args)-i]...)

	return &DomainError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func (e *DomainError) Error() string {
	msg := fmt.Sprintf("code %d %s", e.Code, e.Message)
	if e.Err != nil {
		msg = fmt.Sprintf("%s: %v", msg, e.Err)
	}
	return msg
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

func HasCode(err error, code int) bool {
	var domainError *DomainError
	if errors.As(err, &domainError) {
		return domainError.Code == code
	}
	return false
}
