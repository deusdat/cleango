package cleango

import "testing"

// This test shows how to incorporate repositories and Dependency Injection feature with the use cases.
// In a real project there would be a package called domain. This houses the Entities, Use Cases, and Interfaces that
// define the repositories.
// Repositories define what a use case needs. It does not care how the data is retrieved.
// Another package, maybe http, would implement the repositories. This bifurcation prevents hows from bleeding into
// the whats of the domain.

// ADomainObject is a simple object that presents an entity, which has logic. Very OOP.
type ADomainObject struct {
	StarterValue int
}

func (ado *ADomainObject) Add(newValue int) (int, error) {
	if newValue < 0 {
		return 0, &DomainError{
			Kind:            InvalidInput,
			Message:         "newValue must be zero or more.",
			UnderlyingCause: nil,
			Issues:          nil,
		}
	}
	return ado.StarterValue + newValue, nil
}

type ADomainObjectReader interface {
	Read(id string) (ADomainObject, error)
}

type WithRepoAndEntityInput struct {
	EntityID string
	Amount   int
}

type WithRepoAndEntityUseCase struct {
	repo ADomainObjectReader
}

// Execute
func (s *WithRepoAndEntityUseCase) Execute(input WithRepoAndEntityInput, p Presenter[int]) {
	entity, err := s.repo.Read(input.EntityID)
	if err != nil {
		p.Present(struct {
			Answer int
			Err    error
		}{Answer: 0, Err: err})
		return // This is important. Return to after calling the presenter or else your
		// code will continue. This can lead to presenter errors like writing to a closed
		// response.
	}

	answer, err := entity.Add(input.Amount)
	p.Present(struct {
		Answer int
		Err    error
	}{Answer: answer, Err: err})
}

type DummyRepo struct{}

func (d DummyRepo) Read(id string) (ADomainObject, error) {
	return ADomainObject{
		StarterValue: 10,
	}, nil
}

// The purpose of splitting the use case and presenter into different functions is to show that there may be times
// when you want to reuse a use case and have it present differently based on the context. For example, a REST
// presenter would return JSON, while a website presenter would return a whole HTML page. Keeping them separate from
// the start allows this kind of approach without refactoring.
type SimpleFactory struct {
	// Normally would have singletons like a DB connection, or credentials.
}

func (f *SimpleFactory) UseCase() UseCase[WithRepoAndEntityInput, int] {
	return &WrappingUseCase[WithRepoAndEntityInput, int]{
		Implementation: &WithRepoAndEntityUseCase{
			repo: DummyRepo{},
		},
	}
}

func (f *SimpleFactory) Presenter() Presenter[int] {
	return &presenter{}
}

func TestWithDI(t *testing.T) {
	di := SimpleFactory{}
	useCase := di.UseCase()
	p := di.Presenter()
	useCase.Execute(struct {
		EntityID string
		Amount   int
	}{EntityID: "1000", Amount: 40}, p)

	asp := p.(*presenter)
	if asp.stored != 50 {
		t.Fatalf("amounts did not match expected %d", asp.stored)
	}
}
