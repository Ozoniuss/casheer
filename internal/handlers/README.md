# Testing the handlers

Unfortunately, with this setup testing the handlers is not ideal. Every handler encapsulates the bussiness logic directly, which means that the input comes from an HTTP request and the output is an HTTP response (with whatever gets attached to the context). This makes it really awkward for testing. Of course you can check the context for any attached errors or things like that, but since the handlers ultimately store the bussiness logic, the main concern is to test the bussiness logic thoroughly. This obviously does not happen with this setup.

What you would typically see in a layered architecture is a service layer between the presentation layer (which is handling the HTTP requests) and the bussiness logic layer (which currently is within the handlers). That would be a typical 4-layer architecture, instead of the current standard 3-layer architecture (MVC, more or less). Introducing that layer does require some boilerplate code, but also comes with the advantage that you can test the operations independently, without diving into the HTTP-specific implementation.

Another advantage here would be that the service layer would return service-specific errors, meaning that the tests could look for service-specific errors. It's extremely awkward to check for errors in the context, frankly.

In terms of implementation details, that 4th layer would essentially be an interface exposing all bussiness operations. I'll probably be moving away towards that implementation soon. Testing this really sucks.