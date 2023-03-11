package cleango

import (
	"errors"
	"fmt"
	"testing"
)

func TestPlayground(t *testing.T) {
	de := &DomainError{
		Kind: InvalidInput,
	}

	we := fmt.Errorf("just another wrapping %w", fmt.Errorf("A wrapped error %w", de))
	var couldBe *DomainError
	if errors.As(we, &couldBe) {
		fmt.Printf("error kind was %d\n", couldBe.Kind)
	} else {
		t.Fatal("we should have been the DomainError")
	}
}
