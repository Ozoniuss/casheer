# Spec

Defines the [OpenAPI specification](https://swagger.io/specification/) for
casheer. May be used to generate code for client SDKs, although in Go I
recommend using my own [SDK](../client/) and
[request models](../pkg/casheerapi/).

In production setups APIs are usually versioned, but investing effort to
implement a versioning strategy for this API seemed useless. Adhering to jsonapi
for this personal project was a pain in the ass already.
