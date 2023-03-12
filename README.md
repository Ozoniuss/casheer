# Expense manager

The goal of the application is to manage my own expenses. It basically allows me to plan ahead for each month as well as keep track of the active expenses, by entering them manually, going through my revolut payments etc.

Proposed features
-----------------

There are three main functionalities in the application, that is:

- Make a planning for each month of a year;
- Keep track of a running total for that month by adding expenses;
- Store debts.

The planning is an orientative schedule of future expenses, and helps estimating the amount of money left at the end of each month. This is useful for planning investments, big purchases, vacations etc.

Below is a detailed list of features this application implements.

- [ ] Adding a new entry, assuming all data integrity constraints; 
- [ ] Removing an existing entry. This should soft delete the existing entry, while keeping the transactions that have been recorded for that entry;
- [ ] Listing entries based on filters, including deleted entries;
- [ ] Viewing the expenses for an existing entry, including deleted entries;
- [ ] Updating a single expense (the entry it points to, the value etc.);
- [ ] Adding or removing expenses;
- [ ] Closing an entry (which means that expenses pointing to that entry have to be inserted forcefully);
- [ ] Planning for a month (which is the same as going through the entries of last month interactively and choosing whether to keep or discard them);
- [ ] Closing a month (which means closing all entries of that month);
- [ ] Computing an ongoing total for a month (which is the value of all expenses);

Models
------

A planning is defined for a particular month and year, and has several entries. Each entry holds the following information:

- uuid (used to identify the entry and map expenses to it)
- month
- year
- category
- subcategory
- expected total
- running total
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

Multiple debts for the same person may be issued, and the app will allow to group these debts by person and issue a running total. Depending on what feels best for the user, they might consider updating an ongoing debt if part of it had been settled, or creating a new debt for the same person.