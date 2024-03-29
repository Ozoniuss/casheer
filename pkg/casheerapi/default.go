package casheerapi

type PingResponse struct {
	Info  string    `json:"info"`
	Links PingLinks `json:"links"`
}

type PingLinks struct {
	Entries LinkWithDetails `json:"entries"`
	Debts   LinkWithDetails `json:"debts"`
	Totals  LinkWithDetails `json:"totals"`
}

type LinkWithDetails struct {
	Href    string `json:"href"`
	Details string `json:"details"`
}
