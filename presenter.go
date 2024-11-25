package cleango

// Output is the common result of all use cases.
type Output[T any] struct {
	Answer T
	Err    error
}

// Presenter abstracts the process of taking the answer and presenting it to the caller.
// All use cases should require a presenter as the second argument.
//
// For example, say you are making a website. The use case provides enough information
// for the service to know which Widget it just made. However, you need to go to the database
// to get more information. The presenter implementation can have a database connection, or
// a repository, as part of its initialization. When you construct the presenter, you need to
// include the response object. The presenter will then serialize the response back over the wire.
//
// The result of a presenter is clean, simple code in your web handelers. They simply translate
// the input they get into that which the use case consumers, get the correct presenter from the
// DI engine you use (homespun is fine), and invoke the use case's Execute.
//
// This makes them easy, amazingly easy test. If you wrap your endpoint configurations such that
// you pass the DI as parameter, you can mock/stub it to have a simple implementation that shows
// your code is complete.
//
// Presenters are also responsible for the completion of transactions, either physical or logical.
// In the example above, the presenter could complete the transaction as soon as its invoked.
// The presenter could complete the transaction as part of a "defer". A presenter could complete
// a transaction by completing the db transaction and deleting any temp files.
type Presenter[T any] interface {
	Present(answer Output[T])
}

// PresenterFunc is a simple implementation of Presenter that takes a function as argument. Useful for testing.
type PresenterFunc[T any] struct {
	FN func(Output[T])
}

func (p *PresenterFunc[T]) Present(answer Output[T]) {
	p.FN(answer)
}
