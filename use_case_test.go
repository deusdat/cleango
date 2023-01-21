package go_clean_architecture

import (
	"fmt"
	"strings"
	"testing"
)

type basicUseCase struct {
	shouldPanic bool
}

func (s *basicUseCase) Execute(input string, p Presenter[int]) {
	fmt.Println(strings.ToLower(input))
	if s.shouldPanic {
		panic("implement me")
	}
	p.Present(struct {
		Answer int
		Err    error
	}{Answer: len(input), Err: nil})
}

type presenter struct {
}

func (p *presenter) Present(answer Output[int]) {
	if answer.Err == nil {
		fmt.Printf("That worked fine. Answer is %d\n", answer.Answer)
	} else {
		fmt.Printf("Failed %s\n", answer.Err)
	}
}

func TestWrappingUseCase_Execute(t *testing.T) {
	var wrapAsUseCase UseCase[string, int]
	wrapAsUseCase = &WrappingUseCase[string, int]{
		Implementation: &basicUseCase{},
	}

	wrapAsUseCase.Execute("hello, world", &presenter{})

	wrapAsUseCase = &WrappingUseCase[string, int]{
		Implementation: &basicUseCase{
			shouldPanic: true,
		},
	}

	wrapAsUseCase.Execute("hello, world", &presenter{})
}
