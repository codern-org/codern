package domain

import (
	"errors"
	"fmt"
)

const (
	ErrInternal = 1
	ErrRoute    = 2

	ErrLoggingError     = 1000
	ErrAuthHeader       = 1010
	ErrPayloadParser    = 1011
	ErrPayloadValidator = 1012

	ErrSessionPrefix     = 2000
	ErrSignatureMismatch = 2001
	ErrInvalidSession    = 2002
	ErrSessionExpired    = 2003
	ErrInvalidEmail      = 2010
	ErrDupEmail          = 2011
	ErrUserData          = 2022
	ErrUserPassword      = 2023
)

type GenericError interface {
	Code() int
	error
}

type genericError struct {
	code    int
	message string
	err     error
}

func NewGenericError(code int, message string) GenericError {
	return &genericError{
		code:    code,
		message: message,
		err:     errors.New(message),
	}
}

func (e *genericError) Code() int {
	return e.code
}

func (e *genericError) Error() string {
	return e.err.Error()
}

func (e *genericError) Unwrap() error {
	return e.err
}

type ValidationError interface {
	GenericError
}

type validationError struct {
	namespace string
	field     string
	value     interface{}
}

func NewValidationError(namespace string, field string, value interface{}) ValidationError {
	return &validationError{
		namespace: namespace,
		field:     field,
		value:     value,
	}
}

func (e *validationError) Code() int {
	return ErrPayloadValidator
}

func (e *validationError) Error() string {
	return fmt.Sprintf("Payload validation failed: %s %s %s", e.namespace, e.field, e.value)
}
