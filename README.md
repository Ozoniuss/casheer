# Expense manager

The goal of the application is to manage my own expenses. It basically allows me to plan ahead for each month as well as keep track of the active expenses, by entering them manually, going through my Revolut payments etc.

Proposed features
-----------------

There are three main functionalities in the application, that is:

- Make a planning for each month of a year;
- Keep track of a running total for that month by adding expenses;
- Store debts.

The planning is an orientative schedule of future expenses, and helps estimate the amount of money left at the end of each month. This is useful for planning investments, big purchases, vacations etc.

Below is a detailed list of features this application implements.

Proposed features
-----------------

For entries, I will support:

- [x] Adding a new entry; (1) 
- [x] Removing an existing entry; (2)
- [x] Retrieving a single entry and all its data;
- [x] Updating an entry;
- [x] Listing all entries;
- [x] Filtering entries during listing operations;
- [ ] Including all expenses of an entry in a list operation; 
- [ ] Ordering entries during listing operations;
- [ ] Paginating entries;  


For debts, I will support:

- [x] Adding a new debt; (1) 
- [x] Removing an existing debt; (2)
- [x] Retrieving a single debt and all its data;
- [x] Updating an debt;
- [x] Listing all debts;
- [ ] Ordering debts during listing operations; 
  
For expenses, I will support:

- [x] Adding a new expense, given a valid entry; (1) 
- [x] Removing an existing expense; (2)
- [x] Retrieving a single expense and all its data;
- [x] Updating an expense;
- [x] Listing all expenses for an entry;
- [x] Filtering expenses during listing operations;
- [ ] Ordering expenses during listing operations;

Expenses exist only in the context of an entry, and thus an entry identifier must always be provided, even if it is technically not required to actually perform the operation, such as deleting an expense. This constraint was added for uniformity reasons.

Additional features supported:

- [x] Running total of a period (defined as a month with year). This means, for every valid period, returning a resource containing the expected income and the actual income (based on the `income` category keyword), the expected total of the expenses (not including `income`) and the running total (again, not including `income`). A period is valid if it has at least one entry.

Constraints:

1. All create operations must respect data integrity constraints;
2. All remove operations must be performed via soft deletes (including subsequent cascade operations) in order to minimize losses. A regular cleanup operation will run on the server to delete soft-deleted items from the database.

Models
------

A planning is defined for a particular month and year, and has several entries. Each entry holds the following information:

- uuid (used to identify the entry and map expenses to it)
- month
- year
- category
- subcategory
- expected total
- ~~running total~~ (redundant since it can be computed from expenses, will be removed)
- recurring (if set to true, will be automatically copied to next month)

I don't feel there is a need to have more than a category and a subcategory, otherwise this will turn into a hierarchical mess. The subcategory may also correspond to a one-time event. A good example of an entry that might attract multiple expense would be `transport, fuel` and a good example of an entry that might attract a single expense would be `personal, python-course`. Ideally, the entry category and subcategory should not contain whitespaces to allow for autocomplete in the CLI client.

For each month, there is always a special category called `income`, which should be present for every month and year. There is no restriction on its subcategories, but it almost usually has `salary`. I might consider making additional fixed categories such as `health`, `investments` `food`, but I'm not doing that at the moment.
 
Then we have the expenses, which can be associated with entries. The expenses will contribute to the running total of an entry, either positively or negatively (an expense may have a negative value, as in "a friend paid me back 50$ from months ago"). An expense consists of the following data:

- uuid
- entry uuid
- value
- name
- description (additional information about the expense)
- payment method (may be optional)

Lastly, there's debts. This holds both debts that are owed to me, but also debts that I own to other people. An entry for a debt contains the following data:

- uuid
- person
- amount (positive or negative)
- details

Multiple debts for the same person may be issued, and the app will allow grouping these debts by person and issue a running total. Depending on what feels best for the user, they might consider updating an ongoing debt if part of it had been settled, or creating a new debt for the same person.

Database
--------

The generated database schema can be seen in the image below:

![](assets/database.png)

There's nothing particularly interesting to talk about here except maybe indexing. But before, let's figure out the most common operation I will use the application for, and also estimate how many records the database will hold. 

Starting with the (back-off-the-envelope) estimations:

- Around 5 active debts most of the time. Personally, it makes sense for me to store multiple debts of the same person individually, as it gives me a better view of what is going on. Even considering that, I don't expect to ever have more than, say, 20 active debts;
- A quick search through my revolut history shows me that I have around 100 expenses on average recorded for a month, split accross 30 entries on average. For say 50 years, that is 100 * 12 * 50 = 60000 total records for expenses. This also requires 30 * 12 * 30 = 10,800 records for entries.
- Everything is soft-deleted for about 30 days before being hard-deleted, so it gives me some time to fix mistakes.

Keep in mind that the estimation was for *50 years*. With all these considerations, I don't expect to have more than 100k total records in the database during this lifetime.

It's not even worth talking about queries-per-second. Queries-per-day, maybe. I expect to check the running total from time to time and go through all expenses maybe once or twice a month. During each month, I will be adding all the Revolut expenses one by one, and at the end of the month create the categories for the new month. Thus, the read-to-write ratio isn't really that big, definitely lower than 10. But overall, I would expect somewhere around 1000-2000 queries per month, which is really a laughable load, even if all of them are done in the same day.

So, what are the indexes considerations? The most common operations would be:

- Getting all the expenses for a particular period (year + month), for operations such as "list all expenses" and "get current total";
- Getting details about a particular expense or debt;
- Inserting a new expense.

That being said, I don't think it is even worth considering any indexes at this point. At first, I thought of several strategies, such as indexing (year, month) with a hash index, since many queries will be based on a fixed year or month. Alternatively, indexing them individually with a B+Tree index would help me find expenses in a certain period. But for 100,000 records, I don't really think the additional memory required by these indexes will in the end have a benefit over simply the pages cached by the dbms. For now, the unique index I have acts much more like an additional integrity check for my data.

Now, why not just store all of this in a file? Why use a fully-featured dbms? Well, I do want a relational database since I found a nice way to model my data, but this analysis actually makes me question the choice of Postgres over more lightweight databases such as SQLite, which doesn't even require running a database server. **Backups** and portability are really important feature I want for my database, and SQLite would make those features so much simpler. 

I'll be reading a few reviews of stress tests for SQLite and check out the datatypes and features it offers more in depth, but it's really likely that I will switch from Postgres once this thing is deployed. One of the things that it doesn't seem to support is multiple concurrent writes to the database, which is a feature I will never use. 

TODO: actual storage estimations, potential other features and models for the data.

RESTful service
---------------

Since I'll additionally be presenting this project for a university subject, one of the requirements was to follow the RESTful architectural styles. In short, as defined by Roy Fielding in his [dissertation](https://www.ics.uci.edu/~fielding/pubs/dissertation/fielding_dissertation.pdf), architectural constraints of the following styles should be satisfied:

**The Null Style**

- This is obviously satisfied.

**The Client-Stateless-Server Style**

- This repository contains the backend of the application, whereas a command-line client will be available soon on my GitHub page. The server the system managing and storing expenses, whereas the client will be provided as a means of interaction with the server to perform basic operations and statistics on the available resources.

>

- The stateless property is satisfied by not storing any session information during the interaction between the user and the client. All communication is done via client requests and server responses.

**The Cache Style**
  
- Browsers will automatically implement caching, which can be controlled by the server by specifying HTTP headers such as `Cache-Control`. However, at least in the first iteration, the server will not do any other caching operations, and thus the command-line client will not benefit from caching. But since I plan to use this application only from myself once a month, I'm perfectly fine with not bothering too much about caching.

Perhaps one of the most important feature of a RESTful architecture that is implemented by this backend is providing uniform interfaces for the available resources. Their publicly available representation can be found in the external package [pkg/casheerapi](pkg/casheerapi/), which is going to be imported by the command-line client. The interfaces have been inspired from [jsonapi.org's](https://jsonapi.org/) specification, although many of the points have been dropped in favour of simplicity.

By following these interfaces, a media type in json format had been defined for the resources. Each response will provide state transitions in the form of hyperlinks, which will allow the client to navigate between the available states of the application. This hypertext-enabled behavior is an essential constraint of REST API's, as explained by Roy Fielding himself in [this post](https://roy.gbiv.com/untangled/2008/rest-apis-must-be-hypertext-driven) on his website.

TODO: detail connectors, data elements, views.

Following code principles saves time
--------------------------------------

I was writing the CLI and realized that I wasn't returning errors in the agreed format. This means I had to define an error wrapper in the public package to prefix the error details with the "error" keyword:

```go
type ErrorResponse struct {
	Error Error `json:"error"`
}
```

I realized though I'd dodged a bullet by defining a separate function to render errors to the client:

```go
// EmitError sends an error response back to the client.
func EmitError(ctx *gin.Context, err public.Error) {
	// ctx.JSON(err.Status, gin.H{
	// 	"error": err,
	// })
	ctx.JSON(err.Status, casheerapi.ErrorResponse{
		Error: err,
	})
}
```

I would've had to change every error response had I not factored this out separately.

Obviously a pretty lame example but shows that coding principles are there for a good reason.