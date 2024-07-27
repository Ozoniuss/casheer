# Currency

Defining what a "value" is in this application has made me overthink this more
than I should have. Normally it should be simple: I expect to spend x amount in
my local currency for this entry, and I create expenses for that in the local
currency which I add up to compare with my original expectation. Simple enough.

But even if we approach the problem like this alone, the primary instinct would
be to just store the values as floats, which is actually not recommended. Many
sources like [this one](https://culttt.com/2014/05/28/handle-money-currency-web-applications#use-integers-not-floats)
and [this one](https://cardinalby.github.io/blog/post/best-practices/storing-currency-values-data-types/#bad-choice)
advice against that. TLDR; the operations can be innacurate and you're not 
making use of all the available decimals. I started the app by storing the
smallest unit of my country's RON currency (i.e. a value of 1 RON is stored as
100 in the database).

This was fine at start, but I quickly realised that I was also being paid in 
euros at that time, and some transactions were in euros. Again, the first
instinct is to convert in my country's currency. But then I got all philosophical
into thinking things like "what is a value, in fact" and "what if I have to use
something whose minimum unit might change", and shit like that. I essentially
wanted to generalize the meaning of value -- what if I wanted to make a payment
in oil at some point? 😒...

I didn't really want to study currencies and their minor units to know which
numbers to put in my database. I also didn't really want to store huge numbers
in my database (for example, the smallest Bitcoin unit is 0.00000001 of a BTC)
so I came up with the following system: I would store

- the currency (can be RON, USD, Bitcoin etc.)
- the value
- the exponent (such that value * exponent represents the value in the unit
of that currency)

For example, by storing a value of 150 with currency RON and exponent -2, I'm
actually storing 1.5 RON. This gives me enough flexibility and essentially
allows me to not care at all about the minor unit, or whether if it unexpectedly
changes.

The main purpose of this package is to provide a `Value` struct that encapsulates
such an object defined previously, and provides a set of constructors in order
to always instantiate only valid values of this object. The application will
create values using these constructors, and the database will convert these
values (which are always correct) to its internal representation of the value.

There are some things that I had considered (e.g. the unit of the currency
changing, such as the [revaluation of my country's currency](https://en.wikipedia.org/wiki/Romanian_leu))
or storing the exchange rates at the day of the exchange to be able to accurately,
reproduce and understand the history of my expenses, but I decided to quit
pursuing these thoughts for the sake of my mental health. The real use case of
this application is in the present and short-term past and future, in which case
things like the exact exchange rate at the moment of the transaction doesn't
matter.