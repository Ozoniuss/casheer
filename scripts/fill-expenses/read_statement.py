import csv

import os
import sys
from datetime import datetime

# taken from actual db via an sql query and copy paste, those are just default
# entries. The script doesn't necessarily use one of those.
CATEGORIES = [
    "economies",
    "food",
    "fun",
    "house",
    "income",
    "learning",
    "medic",
    "memberships",
    "misc",
    "personal",
    "transportation",
]
SUBCATEGORIES = [
    "donations",
    "emergency_fund",
    "extra_savings",
    "investment",
    "groceries",
    "takeaway",
    "going_out",
    "electricity",
    "gas",
    "heat",
    "internet",
    "water",
    "salary",
    "scolarship",
    "courses_or_books",
    "meds",
    "backblaze",
    "glovo",
    "google_drive",
    "youtube",
    "idk",
    "teapa",
    "clothing",
    "emag",
    "bolt_bus_or_taxi",
    "fuel",
    "parking",
]


def getMatchesThatStartWith(current: str, wordlist: [str]) -> str:
    matches = []
    for word in wordlist:
        if word.startswith(current):
            matches.append(word)
    return matches


class Expense:
    def __init__(
        self,
        Type: str,
        Product: str,
        StartedDate: str,
        CompletedDate: str,
        Description: str,
        Amount: str,
        Fee: str,
        Currency: str,
        State: str,
        Balance: str,
    ):
        self.Type = Type
        self.Product = Product
        self.Description = Description
        self.Amount = float(Amount)
        self.Fee = float(Fee)
        self.Currency = Currency
        self.State = State

        self.Balance = None
        if Balance != "":
            self.Balance = float(Balance)
        self.StartedDate = None
        self.CompletedDate = None
        if StartedDate != "":
            self.StartedDate = datetime.strptime(StartedDate, "%Y-%m-%d %H:%M:%S")
        if CompletedDate != "":
            self.CompletedDate = datetime.strptime(CompletedDate, "%Y-%m-%d %H:%M:%S")

    def __str__(self):
        return f"({str(self.StartedDate.date()) :<10}) {self.Currency :<4} {int(self.Amount * 100) :<7} {int(self.Fee*100) :<2} {'[' + self.Description + ']' :<25}"


filename = os.getenv("EXPENSES_STATEMENT")


def askInput(amount: int, currency: str, payment_method: str, name: str):
    keep = input("keep? >")
    if keep == "n":
        raise Exception("done")
    actual_amount = input(f"amount? (default {amount}) >")
    if actual_amount == "":
        actual_amount = amount
    actual_amount = int(actual_amount)
    actual_currency = input(f"currency? (default {currency}) >")
    if actual_currency == "":
        actual_currency = currency
    actual_payment_method = input(f"payment method? (default {payment_method}) >")
    if actual_payment_method == "":
        actual_payment_method = payment_method
    category = input("category? >")
    if category not in CATEGORIES:
        category = input(
            f"are you sure? you typed {category}, possible matches {getMatchesThatStartWith(category, CATEGORIES)}\nretype if sure >"
        )
        if category not in CATEGORIES:
            print("no match, adding category")
            CATEGORIES.append(category)
    subcategory = input("subcategory? >")
    if subcategory not in SUBCATEGORIES:
        subcategory = input(
            f"are you sure? you typed possible matches {getMatchesThatStartWith(subcategory, SUBCATEGORIES)}\nretype if sure >"
        )
        if subcategory not in SUBCATEGORIES:
            print("no match, adding subcategory")
            SUBCATEGORIES.append(subcategory)
    actual_name = input(f"name? (default {name})>")
    if actual_name == "":
        actual_name = name
    description = input("description? >")

    return (
        category,
        subcategory,
        actual_name,
        description,
        actual_amount,
        actual_currency,
        actual_payment_method,
    )


with open(filename, "r") as csvfile:
    spamreader = csv.reader(csvfile)
    expenses: list[Expense] = []
    for idx, row in enumerate(spamreader):
        # ignore header
        if idx == 0:
            continue
        # vault transfers or card topups, ignore those
        if (
            row[4] == "To RON Bani de rontanele"
            or row[4] == "To RON"
            or row[4].startswith("Top-Up by *")
        ):
            continue
        exp = Expense(*row)
        expenses.append(exp)
    expenses.sort(key=lambda expense: expense.StartedDate)

with open("out.txt", "w") as f:
    for e in expenses:
        print(e)
        try:
            cat, subcat, name, desc, act_am, act_cur, act_pymt = askInput(
                -int(e.Amount * 100), e.Currency, "card", e.Description
            )
            f.write(f"{cat},{subcat},{name},{desc},{act_am},{act_cur},{act_pymt}\n")
        except Exception as e:
            if str(e) == "done":
                print("skipped")
