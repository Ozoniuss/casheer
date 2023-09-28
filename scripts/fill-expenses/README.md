# Fill expenses

While this app is still in development but the database was already in use, I
wrote a couple of short scripts that helped me introduce the expenses in the
database with the help of the Go client and the revolut statement.

### [Create default expenses Go script](create_default_expenses.go)

This script introduces the entries of this month in the database based on the
default entries json template. I've obviously not added my personal template,
but see how a sample similar one looks [here.](../../sample_entry_template.json)

### [Add to db Go script](add_to_db.go)

This script achieved an important personal milestone for me, and that is, adding
the very first expenses in the production database! I used these scripts to go
through my revolut statement and add the expenses interactively one by one.
There are obviously many more improvements I could do to this, but I just wanted
a proof of concept and also do some analysis on my expenses this month. Now that
I know it works, I've set the building blocks of using the Go client and will
start working on its actual use case: the CLI.

Not only that, but this also helped me do a full end-to-end test by validating
the expenses that were being continuously created by the script agains my real
Revolut account, and I was really glad to see it didn't miss a single one.

### [Read statement Python script](read_statement.py)

Back when I was experimenting with extracting my Revolut expenses, I was set on
Python and thus I wrote a Python script that scraped my Revolut statement. This
script is used to:

- filter all the relevant expenses from my statement;
- go through each one and allow to interactively fill my expense details;
- either do a command-line autocomplete or an automated spell mistake fix for
  the category and subcategory.

Once all expenses are parsed, a file will be generated with the expenses in a
good format, and another Go script (which doesn't exist it yet) will add them to
the database with the help of the client.
