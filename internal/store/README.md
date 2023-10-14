# Store

In theory, this package is responsible for managing database interaction. It
should provide a `store` abstraction which handles connecting to an external
storage like a database and defines various operations which can be performed
against that storage that are related to the application: reads, writes, queries
etc. The idea is that the bussiness layer would not depend directly on the
persistence layer, but on the `store` abstraction, allowing for:

- easily swapping storages by implementing the `store` abstraction;
- defining some common errors that the bussiness layer needs to be aware of
  which wrap the errors thrown by the underlying storage libraries, leading to
  effective error handling at the upper layer by only treating the ones exposed
  by the abstraction;
- defining `store` mocks to allow for true unit tests for the bussiness logic
  that depends on external storage.

That is, in theory. In practice, I've only included here some helpers for
connecting to the database and running some migrations. When I started the
project I didn't think too much of the storage layer and thought I'd only ever
use a postgres database, and that implementing this abstraction would be a waste
of time. At some point I switched to SQLite and I was lucky enough that `gorm`
abstracted enough of the database interaction that I only needed to change the
driver. But, this is still not ideal because I'm still handling `gorm` errors
directly in the bussiness layer, I'm testing with a real database (which is
reasonable since SQLite only uses a file but doesn't scale well to other
databases) and I have to write actual queries in the handlers. I should probably
move away from this, but the code is already running for my production expenses
database and refactoring doesn't seem easy so yeah.
