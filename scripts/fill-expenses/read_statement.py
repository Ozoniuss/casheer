import csv
import os
from datetime import datetime

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
        if Balance != '':
            self.Balance = float(Balance)
        self.StartedDate = None
        self.CompletedDate = None
        if StartedDate != '':
            self.StartedDate = datetime.strptime(StartedDate, '%Y-%m-%d %H:%M:%S')
        if CompletedDate != '':
            self.CompletedDate = datetime.strptime(CompletedDate, '%Y-%m-%d %H:%M:%S')

    def __str__(self):
        return f"{self.Type} {self.Product} {self.Currency} {int(self.Amount * 100)} {int(self.Fee*100)} [{self.Description}]"

filename = os.environ("EXPENSES_STATEMENT")
card_ending_no = os.environ("CARD_ENDING_NO")

with open(filename) as csvfile:
    spamreader = csv.reader(csvfile)
    expenses: list[Expense] = []
    for idx, row in enumerate(spamreader):
        if idx == 0:
            continue
        if row[4] == 'To RON Bani de rontanele' or row[4] == 'To RON' or row[4] == f'Top-Up by *{card_ending_no}':
            continue
        exp = Expense(*row)
        expenses.append(exp)
        print(exp)
    