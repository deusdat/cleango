# cleango

The purpose of this library is to provide a set of pre-build artifacts that
conform to the key elements of Clean Architecture as described by Uncle Bob
[here.](https://blog.cleancoder.
com/uncle-bob/2012/08/13/the-clean-architecture.html) If you are unfamiliar 
with the concept please see this presentation. 

[![IMAGE ALT TEXT HERE](https://img.youtube.com/vi/o_TH-Y78tt4/0.jpg)](https://www.youtube.com/watch?v=o_TH-Y78tt4)


Many implementations that claim to pick up the mantle lack the ideas of a 
Presenter, Use Case and inversion of control. They instead have a fairly 
procedural approach where the controller instantiates an object or calls a 
function where either are directly invoked, waits for a response, handles 
errors (maybe a middleware does that), creates a new response and sends it 
over the wire. The result is a hard-to-test set of assumptions and coupling. 

Clean, as its name implies, works to stratify your code into isolated chunks 
that coordinate via **Dependency Inject**. A controller gets a dependency 
injector passed to it to create Use Cases and Presenters. Those in turn have 
parts injected into them like a Repository to get data. Those Repositories 
are defined by interfaces. The Use Case doesn't care that the Repository 
interacts with a Database, nor does it care if it's calling a ReST service. 
The domain defines the interface of what it needs, and something, somewhere 
implements that interface, is properly constructed and injected into the Use 
Case.

Behind the scenes you can get fancy with the DI system. You can use [Unit of 
Work](https://www.oreilly.com/library/view/beginning-solid-principles/9781484218488/A416860_1_En_10_Chapter.html) 
to coordinate multiple changes like writing to more than one collection in 
an [ArangoDB](https://www.arangodb.com) or completing a database transaction 
and deleting temp files.
