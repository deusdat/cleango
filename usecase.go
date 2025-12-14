package cleango

import (
	"context"
	"fmt"
)

// UseCase defines the core of the computational space. It only works on input and returns an Answer to the presenter
// via dependency inject. Notice this method does not return a value. It doesn't even return an error. Any errors
// encountered during calculation should be sent to the presenter via the Present method.
type UseCase[Input any, Answer any] interface {
	Execute(Input, Presenter[Answer])
}

// Same idea as the regular use case but context aware.
type UseCaseWithContext[Input any, Answer any] interface {
	Execute(context.Context, Input, PresenterWithContext[Answer])
}

// WrappingUseCase a use case that holds another use case in order to make sure presenter is always caused with
// an error in case of a panic.
type WrappingUseCase[Input any, Answer any] struct {
	Implementation UseCase[Input, Answer]
}

var RecoveryMessage = "wrapper recovery"

func (w *WrappingUseCase[Input, Answer]) Execute(input Input, p Presenter[Answer]) {
	defer func() {
		if r := recover(); r != nil {
			var blank Answer
			p.Present(struct {
				Answer Answer
				Err    error
			}{
				Answer: blank,
				Err: &DomainError{
					Kind:            System,
					Message:         RecoveryMessage,
					UnderlyingCause: fmt.Errorf("%s", r),
					Issues:          nil,
				},
			})
		}
	}()
	w.Implementation.Execute(input, p)
}

// FunctionalUseCase is a generic implementation of UseCase that takes a function to act as the Execute method.
type FunctionalUseCase[Input any, Answer any] struct {
	ExecuteFunc func(Input) (Answer, error)
}

// Execute runs the provided function and passes its result to the presenter.
func (f *FunctionalUseCase[Input, Answer]) Execute(input Input, p Presenter[Answer]) {
	answer, err := f.ExecuteFunc(input)
	p.Present(Output[Answer]{
		Answer: answer,
		Err:    err,
	})
}

// NewFunctionalUseCase creates a new FunctionalUseCase with the provided function.
func NewFunctionalUseCase[Input any, Answer any](executeFunc func(Input) (Answer, error)) *FunctionalUseCase[Input, Answer] {
	return &FunctionalUseCase[Input, Answer]{
		ExecuteFunc: executeFunc,
	}
}
