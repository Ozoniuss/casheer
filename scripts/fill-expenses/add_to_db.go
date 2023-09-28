package main

// don't expect much from this, it's 12 am and I just quickly wrote this
// for the sake of getting them in the db

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Ozoniuss/casheer/client/httpclient"
)

// what i care about
type data struct {
	category       string
	subcategory    string
	name           string
	description    string
	amount         int
	currency       string
	payment_method string
}

// too lazy to make this nice, just get data from line and insert in db with
// client
func insertLine(client *httpclient.CasheerHTTPClient, line string) {
	fields := strings.Split(line, ",")
	amnt, err := strconv.Atoi(fields[4])
	if err != nil {
		panic(err)
	}
	expdata := data{
		category:       fields[0],
		subcategory:    fields[1],
		name:           fields[2],
		description:    fields[3],
		amount:         int(amnt),
		currency:       fields[5],
		payment_method: fields[6],
	}
	_, err = client.CreateBasicExpenseWithoutId(
		expdata.category,
		expdata.subcategory,
		9,
		2023,
		expdata.name,
		expdata.description,
		expdata.payment_method,
		expdata.amount,
		expdata.currency,
	)
	if err != nil {
		// probably entry doesnt exist, will add a method to automatically
		// create it
		fmt.Printf("could not create expense: %s\n", err.Error())
		_, err := client.CreateEntry(9, 2023, expdata.category, expdata.subcategory, expdata.amount, false)
		if err != nil {
			fmt.Println("not good")
			panic(err)
		}
		_, err = client.CreateBasicExpenseWithoutId(
			expdata.category,
			expdata.subcategory,
			9,
			2023,
			expdata.name,
			expdata.description,
			expdata.payment_method,
			expdata.amount,
			expdata.currency,
		)
		if err != nil {
			fmt.Println("not good again")
			panic(err)
		}
	}
}

func main() {
	f, _ := os.Open("out_final.txt")
	defer f.Close()

	cclient, err := httpclient.NewCasheerHTTPClient(
		httpclient.WithCustomAuthority("localhost:8033"),
	)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		insertLine(cclient, line)
	}

}
