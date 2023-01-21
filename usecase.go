package go_clean_architecture

// UseCase defines the core of the computational space. It only works on input and returns an Answer to the presenter
// via dependency inject. Notice this method does not return a value. It doesn't even return an error. Any errors
// encountered during calculation should be sent to the presenter via the Present method.
type UseCase[Input any, Answer any] interface {
	Execute(Input, Presenter[Answer])
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
					UnderlyingCause: r,
					Issues:          nil,
				},
			})
		}
	}()
	w.Implementation.Execute(input, p)
}
