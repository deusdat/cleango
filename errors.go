package cleango

import (
	"errors"
	"fmt"
	"strings"
)

type ErrorKind int

const (
	// System means that something happened that prevents the code continuing.
	System ErrorKind = iota
	// NotFound indicates the system could not find an expected resource
	NotFound
	// InvalidInput indicates that the input was expected in some way
	InvalidInput
	// Duplicate indicates that a resource already exists.
	Duplicate
)

// ValidationIssue detail about where an input is wrong. Recommended to use a json path.
type ValidationIssue struct {
	Path string
	// Message basically anything you want. You can put codes for i18n lookups.
	Message string
	// Min allows you to specify the lower bound of field.
	Min int
	// Max allows you to specify the upper bound of a field.
	Max int
}

// DomainError is the only error definition that is used by domain layer. All repositories and use cases should
// create these errors are necessary. For the most part, use cases probably won't need to interrogate errors.
// Make sure you properly wrap the errors and use errors.Is or errors.As to get the specific details when appropriate.
type DomainError struct {
	// Kind which specific error is at hand.
	Kind ErrorKind
	// Message is a human-readable message describing the cause of the error.
	Message string
	// UnderlyingCause the source error that caused this. Not to be confused with a wrapped error. This is error is
	// optional and to be used, if necessary, to provide an outer layer detailed information that might not need to
	// be communicated with the caller.
	UnderlyingCause error
	// Issues an issues that occurred while validating input. Should be paired with InvalidInput, but it's your
	// code base.
	Issues []ValidationIssue
}

var toHuman map[ErrorKind]string = make(map[ErrorKind]string)

// InvalidInputKindAsString translates the Kind iota into a human-readable. Update for non-English.
var InvalidInputKindAsString = "invalid input"

// SystemKindAsString translates the Kind iota into human-readable. Update for non-English
var SystemKindAsString = "system"

// NotFoundKindAsString translates the Kind iota into human-readable. Update for non-English
var NotFoundKindAsString = "not found"

func init() {
	toHuman[InvalidInput] = InvalidInputKindAsString
	toHuman[System] = SystemKindAsString
	toHuman[NotFound] = NotFoundKindAsString
}
func (d *DomainError) Error() string {
	if d == nil {
		return "domain error was nil"
	}
	kind := toHuman[d.Kind]
	errMsg := "[%s - %s]"
	all := []any{kind, d.Message}
	if d.UnderlyingCause != nil {
		all = append(all, d.UnderlyingCause.Error())
		errMsg = "[%s - %s (%s)]"
	}

	return fmt.Sprintf(errMsg, all...)
}

var ToDomainErrorMessage = "converted error"

// ToDomainError will wrap an error. If the error is not a domain error,
// it will create one with the underlying cause set to the original err value.
// This way all errors will unwrap to a DomainError.
//
// If err is nil, return nil. Supports cases where Anwer and Err are set in presenter
// after a simple service invocation.
func ToDomainError(extraMessage string, err error) error {
	if err == nil {
		return nil
	}

	var possibleDomainError *DomainError
	if errors.As(err, &possibleDomainError) {
		if !strings.Contains(extraMessage, "%w") {
			// Make sure the error is properly wrapped.
			extraMessage += " (%w)"
		}
		return fmt.Errorf(extraMessage, err)
	}
	return &DomainError{
		Kind:            System,
		Message:         extraMessage,
		UnderlyingCause: err,
		Issues:          nil,
	}
}
