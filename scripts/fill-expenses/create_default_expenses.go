package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Ozoniuss/casheer/client/httpclient"
)

type ExpData struct {
	ExpectedTotal float32 `json:"expected_total"`
	Recurring     bool    `json:"recurring"`
	Currency      string  `json:"currency"`
}

func run() error {
	f, err := os.Open("../../default.json")
	if err != nil {
		return fmt.Errorf("could not open file: %s", err.Error())
	}
	defer f.Close()

	var data = make(map[string]map[string]ExpData)

	dec := json.NewDecoder(f)
	err = dec.Decode(&data)
	if err != nil {
		return fmt.Errorf("could not decode: %s", err.Error())
	}

	cclient, err := httpclient.NewCasheerHTTPClient(
		httpclient.WithCustomAuthority("localhost:8033"),
	)
	if err != nil {
		return fmt.Errorf("could not create client: %s", err.Error())
	}

	for category, subcategories := range data {
		for subcategory, data := range subcategories {
			if data.Currency == "" {
				data.Currency = "RON"
			}
			_, err := cclient.CreateEntry(10, 2024, category, subcategory, int(data.ExpectedTotal*100), data.Currency, data.Recurring)
			if err != nil {
				fmt.Printf("could not create entry: %s\n", err.Error())
			}
		}
	}
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
