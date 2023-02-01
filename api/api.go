package api

import (
	"net/http"

	"github.com/dghubble/sling"
)

type Hash map[string]interface{}

type API struct {
	client *sling.Sling
}

type APIError struct {
	Message  string              `json:"error"`
	Messages []map[string]string `json:"errors"`
}

func (e APIError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	for _, m := range e.Messages {
		for _, v := range m {
			return v
		}
	}
	return "unknown error response"
}

func New(customerBase, customerApiKey string) *API {
	sling := sling.New().
		Client(http.DefaultClient).
		Base(customerBase).
		SetBasicAuth("", customerApiKey).
		Set("User-Agent", "terraform")

	return &API{
		client: sling,
	}
}
