package main

import (
	"encoding/json"
	"os"

	"github.com/Ozoniuss/casheer/client/httpclient"
)

type ExpData struct {
	ExpectedTotal float32 `json:"expected_total"`
	Recurring     bool    `json:"recurring"`
}

func main() {
	f, _ := os.Open("default.json")
	defer f.Close()

	var data = make(map[string]map[string]ExpData)

	dec := json.NewDecoder(f)
	err := dec.Decode(&data)
	if err != nil {
		panic(err)
	}

	cclient, err := httpclient.NewCasheerHTTPClient(
		httpclient.WithCustomAuthority("localhost:8033"),
	)
	if err != nil {
		panic(err)
	}

	for category, subcategories := range data {
		for subcategory, data := range subcategories {
			_, err := cclient.CreateEntry(9, 2023, category, subcategory, int(data.ExpectedTotal*100), data.Recurring)
			if err != nil {
				panic(err)
			}
		}
	}

}
