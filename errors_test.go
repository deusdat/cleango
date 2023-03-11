package cleango

import (
	"errors"
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

	de := ToDomainError(
		"wrapping another domain",
		&DomainError{
			Kind:    InvalidInput,
			Message: "bad param {jimmy}",
		})
	if errors.As(de, &asDomain) {
		if !errors.As(asDomain.UnderlyingCause, &asDomain) {
			t.Fatal("should have been domain")
		}
		if asDomain.Kind != InvalidInput {
			t.Fatal("incorrect mapping")
		}
	} else {
		t.Fatalf("could not convert domain error")
	}
}
