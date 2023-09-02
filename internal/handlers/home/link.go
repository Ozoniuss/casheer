package home

import (
	"net/url"

	"github.com/Ozoniuss/casheer/pkg/casheerapi"
)

func NewHomeLink(baseURL *url.URL) casheerapi.HomeLink {
	return casheerapi.HomeLink{
		Href:  baseURL.JoinPath("home").String(),
		Title: "Home page of casheer API.",
	}
}
