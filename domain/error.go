package domain

import (
	"errors"
	"fmt"
)

const (
	ErrInternal = 1
	ErrRoute    = 2
	ErrParam    = 3

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

	ErrWorkspaceNotFound = 3000
)

type DomainError interface {
	Code() int
	error
}

type domainError struct {
	code    int
	message string
	err     error
}

func NewError(code int, message string) DomainError {
	return &domainError{
		code:    code,
		message: message,
		err:     errors.New(message),
	}
}

func (e *domainError) Code() int {
	return e.code
}

func (e *domainError) Error() string {
	return e.err.Error()
}

func (e *domainError) Unwrap() error {
	return e.err
}

func HasErrorCode(err error, code int) bool {
	var domainError DomainError
	if errors.As(err, &domainError) {
		return domainError.Code() == code
	}
	return false
}

type ValidationError interface {
	DomainError
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
