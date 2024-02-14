package cleango

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrors(t *testing.T) {
	e := ToDomainError("no wrap message", errors.New("not domain"))
	var asDomain *DomainError
	if !errors.As(e, &asDomain) {
		t.Fatal("should have created new domain error")
	}
	if asDomain.UnderlyingCause == nil {
		t.Fatal("underlying cause was not preserved")
	}
	if asDomain.UnderlyingCause.Error() != "not domain" {
		t.Fatalf("unknown underlying cause %s", asDomain.UnderlyingCause)
	}

	if asDomain.Error() != "[system - converted error (not domain)]" {
		t.Fatal("did not nest call stack properly", asDomain.Error())
	}

	de := ToDomainError(
		"wrapping another domain",
		&DomainError{
			Kind:    InvalidInput,
			Message: "bad param {jimmy}",
		})
	if !errors.As(de, &asDomain) ||
		!errors.Is(asDomain, errors.Unwrap(de)) {
		t.Fatalf("unwrapped message did not match source")
	}
}

func TestDeepIssues(t *testing.T) {
	dbLikeErr := fmt.Errorf("failed to connect to data source")
	wrapper1 := ToDomainError("converted to domain error", dbLikeErr)
	useCaseErrWrapper := ToDomainError("wrapper at the use case level", wrapper1)
	if useCaseErrWrapper.Error() != "wrapper at the use case level ([system - converted to domain error (failed to connect to data source)])" {
		t.Fatal("message didn't wrap properly | ", useCaseErrWrapper.Error())
	}
}
