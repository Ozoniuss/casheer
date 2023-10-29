# Golang Casheer Client

At work, we call this "SDK driven development". And it's fucking awesome. The
idea is pretty simple: you provide a client for the API, which is essentially
an interface of methods that can have one or multiple implementations. Then
people can download your SDK to interact directly with your API from their code.
Many products [do this](https://aws.amazon.com/sdk-for-go/), and for good reasons:

- The interface offers an opinionated way to interact with your application 
using the programming language the SDK is written in, abstracting away things 
like the HTTP client or building your own requests and handling responses. If
you're using the SDK, upgrading its version is enough -- you no longer need to
understand the API specs

- Therefore, it improves productivity by simplifying the interaction, and is 
also less error-prone, since the API developers do the heavy-lifting by 
implementing and maintaining this client. They know and understand their API
best, after all.

- It also benefits the API developers by simplifying testing **significanlty**.
They can use the SDK to write full integration tests, by calling the API with
the help of the SDK.

- It allows for multiple implementations of the same interface by the
developers, allowing you to choose the implementation that suits best your use
case. A very powerful use case is defining a "mock" implementation which doesn't
really communicate with any backend. This enables true unit tests at the
bussiness logic level, if that bussiness logic involves communicating with 
this API. For example, if your handlers involve at some point calling this API,
and you want to unit-test your handlers, you can use a mock implementation so
you don't have to set up any other dependencies. At this stage, you can assume
that the client is tested and works, there is no need for a full integration
test.

- Not only that, but the mock implementation can also be used to ensure that 
the return types of this API interface are what they're intended to be.
Regardless of the implementation you're using, if starting with the same data,
doing the same operations should result in the same outcome .A nice
pattern is to set up so-called "contract tests", where you test all
implementations of the client at the same time, by just calling the interface
methods in your tests. If the "real" client is part of the tests, this also
proves that any "mock" implementation works as expected and can be trusted when
unit testing the upper layer, e.g. the bussiness logic.

- The contract tests also play the role of documentation: they are a great place
to look at in order to understand how to use the client.

Considerations specific to Casheer
----------------------------------

This client is in development and at the moment doesn't have a fake
implementation, but that will be coming soon. Though, you can still check out
[the end-to-end tests](../e2e/) which showcase how the points made above about
testing have been applied to this project.

Also, I was lasy and the interface methods return the same structs that I
created to model the jsonapi responses. I don't really like that and I think it
would be better practice if the client defined its own models and errors to 
match the intended use cases. The client users should not be concerned of the
API representation of the data ([unifrom interface](https://stackoverflow.com/questions/25172600/rest-what-exactly-is-meant-by-uniform-interface) in REST slang), such as knowing
that it follows jsonapi, which requires understanding that standard to know
where to take the fields from.