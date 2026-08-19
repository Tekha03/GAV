package apperrors

import (
	"context"
	"errors"
)

type Error struct {
	Code     Code
	Category Category
	Message  string
	Details  map[string]any
	Cause    error
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	return e.Message
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Cause
}

func (e *Error) Is(target error) bool {
	if e == nil {
		return false
	}

	targetError, ok := target.(*Error)
	return ok && targetError != nil && e.Code == targetError.Code
}

type Option func(*Error)

func WithDetail(key string, value any) Option {
	return func(e *Error) {
		if e.Details == nil {
			e.Details = make(map[string]any)
		}

		e.Details[key] = value
	}
}

func WithDetails(details map[string]any) Option {
	return func(e *Error) {
		if len(details) == 0 {
			return
		}

		if e.Details == nil {
			e.Details = make(map[string]any, len(details))
		}

		for key, value := range details {
			e.Details[key] = value
		}
	}
}

func New(def Definition, message string, options ...Option) *Error {
	e := &Error{
		Code:     def.Code,
		Category: def.Category,
		Message:  message,
	}

	applyOptions(e, options)
	return e
}

func Wrap(def Definition, message string, cause error, options ...Option) error {
	if cause == nil {
		return nil
	}

	e := &Error{
		Code:     def.Code,
		Category: def.Category,
		Message:  message,
		Cause:    cause,
	}

	applyOptions(e, options)
	return e
}

func CodeOf(err error) (Code, bool) {
	appError, ok := asError(err)
	if !ok {
		return "", false
	}

	return appError.Code, true
}

func CategoryOf(err error) (Category, bool) {
	appError, ok := asError(err)
	if !ok {
		return "", false
	}

	return appError.Category, true
}

func DetailsOf(err error) (map[string]any, bool) {
	appError, ok := asError(err)
	if !ok {
		return nil, false
	}

	if len(appError.Details) == 0 {
		return nil, true
	}

	details := make(map[string]any, len(appError.Details))
	for key, value := range appError.Details {
		details[key] = value
	}

	return details, true
}

func Normalize(err error) *Error {
	if err == nil {
		return nil
	}

	if appError, ok := asError(err); ok {
		return appError
	}

	switch {
	case errors.Is(err, context.Canceled):
		return newWithCause(RequestCancelled, "request cancelled", err)
	case errors.Is(err, context.DeadlineExceeded):
		return newWithCause(RequestTimeout, "request timeout", err)
	default:
		return newWithCause(Internal, "internal server error", err)
	}
}

func asError(err error) (*Error, bool) {
	if err == nil {
		return nil, false
	}

	var appError *Error
	if !errors.As(err, &appError) || appError == nil {
		return nil, false
	}

	return appError, true
}

func applyOptions(e *Error, options []Option) {
	for _, option := range options {
		if option != nil {
			option(e)
		}
	}
}

func newWithCause(def Definition, message string, cause error) *Error {
	return &Error{
		Code:     def.Code,
		Category: def.Category,
		Message:  message,
		Cause:    cause,
	}
}
