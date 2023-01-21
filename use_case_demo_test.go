package cleango

import (
	"fmt"
	"strings"
	"testing"
)

type basicUseCase struct {
	shouldPanic bool
}

// Execute a simple use case example that shows the flow of input, business logic, and presenter invocation.
func (s *basicUseCase) Execute(input string, p Presenter[int]) {
	fmt.Println(strings.ToLower(input))
	// You'd never want to do this. However, it shows the WrapperUseCase in action.
	if s.shouldPanic {
		panic("implement me")
	}
	p.Present(struct {
		Answer int
		Err    error
	}{Answer: len(input), Err: nil})
}

type presenter struct {
	errOccurred bool
}

func (p *presenter) Present(answer Output[int]) {
	if answer.Err == nil {
		fmt.Printf("That worked fine. Answer is %d\n", answer.Answer)
	} else {
		fmt.Printf("Failed %s\n", answer.Err)
		p.errOccurred = true
	}
}

func TestWrappingUseCase_Execute(t *testing.T) {
	var wrapAsUseCase UseCase[string, int]
	wrapAsUseCase = &WrappingUseCase[string, int]{
		Implementation: &basicUseCase{},
	}

	p := &presenter{}
	wrapAsUseCase.Execute("hello, world", p)
	if p.errOccurred {
		t.Fatal("should not have failed")
	}

	wrapAsUseCase = &WrappingUseCase[string, int]{
		Implementation: &basicUseCase{
			shouldPanic: true,
		},
	}

	wrapAsUseCase.Execute("hello, world", p)
	if !p.errOccurred {
		t.Fatal("should have failed")
	}
}
