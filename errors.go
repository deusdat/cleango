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
)

// ValidationIssue detail about where an input is wrong. Recommended to use a json path.
type ValidationIssue struct {
	Path string
	// Message basically anything you want. You can put codes for i18n lookups.
	Message string
	// Cardinality allows to contextualize messages like, too many items. No more than Cardinality.
	Cardinality int
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
	return fmt.Sprintf("%s - %s", kind, d.Message)
}

var ToDomainErrorMessage = "converted error"

// ToDomainError will wrap an error. If the error is not a domain error, it will create one with the
// underlying cause set to the original err value.
func ToDomainError(extraMessage string, err error) error {
	var possibleDomainError *DomainError
	if !strings.Contains(extraMessage, "%w") {
		// Make sure the error is properly wrapped.
		extraMessage += extraMessage + " %w"
	}
	if errors.As(err, &possibleDomainError) {
		return fmt.Errorf(extraMessage, err)
	}
	return fmt.Errorf(extraMessage, &DomainError{
		Kind:            System,
		Message:         ToDomainErrorMessage,
		UnderlyingCause: err,
		Issues:          nil,
	})
}
