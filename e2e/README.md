# E2E

End to end tests. Since it's so simple to run casheer (the database is just a 
file and doesn't require a separate process), putting these tests up was not 
complicated. They simply run the app on some port and use the [client](../client/)
to communicate with the backend, checking that the response is the expected one.