// Package tools provides MCP tool handlers for Transmission operations.
package tools

import "github.com/cockroachdb/errors"

// ErrValidation indicates invalid parameters provided by the caller.
var ErrValidation = errors.New("validation error")

// ErrNegativeLimit is returned when a numeric limit parameter is negative.
var ErrNegativeLimit = errors.New("numeric limits must not be negative")

// ErrAbsolutePathRequired is returned when a path parameter is not absolute.
var ErrAbsolutePathRequired = errors.New("path must be absolute (start with /)")

// ErrTransmission indicates a failure communicating with the Transmission RPC API.
var ErrTransmission = errors.New("transmission request error")

// validationErr marks an error as a validation error.
func validationErr(err error) error {
	//nolint:wrapcheck // Mark adds a sentinel category, the caller already provides context.
	return errors.Mark(err, ErrValidation)
}

// transmissionErr wraps a message and underlying error as a Transmission request error.
func transmissionErr(msg string, err error) error {
	//nolint:wrapcheck // Mark adds a sentinel category on top of Wrap which provides context.
	return errors.Mark(errors.Wrap(err, msg), ErrTransmission)
}
