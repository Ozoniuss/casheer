package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/Ozoniuss/casheer/client/httpclient"
)

func askWithPrompt(prompt string, s *bufio.Scanner) string {
	fmt.Print(prompt)
	s.Scan()
	return s.Text()
}

type debt struct {
	name     string
	value    int
	currency string
	details  string
}

type expense struct {
	category       string
	subcategory    string
	name           string
	description    string
	amount         int
	currency       string
	payment_method string
}

func askForDebt(scan *bufio.Scanner) debt {

	name := askWithPrompt("Who's the debt for? ", scan)
	value := askWithPrompt("What's the debt value? ", scan)
	currency := askWithPrompt("What's the currency? ", scan)
	details := askWithPrompt("Any details? ", scan)

	val, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}

	return debt{name: name, value: val, currency: currency, details: details}
}

func insertDebt(cclient *httpclient.CasheerHTTPClient, scan *bufio.Scanner) {
	debt := askForDebt(scan)
	resp, err := cclient.CreateDebt(debt.name, debt.details, debt.value, debt.currency, -2)
	fmt.Println(resp, err)
}

func askForExpense(scan *bufio.Scanner) expense {
	category := askWithPrompt("category: ", scan)
	subcategory := askWithPrompt("subcategory: ", scan)
	name := askWithPrompt("name: ", scan)
	description := askWithPrompt("description: ", scan)
	amount := askWithPrompt("amount (in minor unit): ", scan)
	currency := askWithPrompt("currency: ", scan)
	payment_method := askWithPrompt("payment method: ", scan)

	amt, err := strconv.Atoi(amount)
	if err != nil {
		panic(err)
	}
	return expense{
		category:       category,
		subcategory:    subcategory,
		name:           name,
		description:    description,
		amount:         amt,
		currency:       currency,
		payment_method: payment_method,
	}
}

func insertExpense(cclient *httpclient.CasheerHTTPClient, scan *bufio.Scanner) {
	expense := askForExpense(scan)
	resp, err := cclient.CreateBasicExpenseWithoutId(
		expense.category,
		expense.subcategory,
		11,
		2023,
		expense.name,
		expense.description,
		expense.payment_method,
		expense.amount,
		expense.currency,
	)
	fmt.Println(resp, err)
}

func main() {
	cclient, err := httpclient.NewCasheerHTTPClient(
		httpclient.WithCustomAuthority("localhost:8033"),
	)
	if err != nil {
		panic(err)
	}

	scan := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("what do you want to insert? [d,ex,q] ")
		scan.Scan()
		if scan.Text() == "d" {
			insertDebt(cclient, scan)

		} else if scan.Text() == "q" {
			return
		} else if scan.Text() == "ex" {
			insertExpense(cclient, scan)
		}
	}

}
