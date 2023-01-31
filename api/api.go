package api

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

type Json map[string]interface{}

type API struct {
	client *sling.Sling
}

type APIError struct {
	Message string `json:"error"`
}

func (e APIError) Error() string {
	return e.Message
}

type APIErrors struct {
	Messages []map[string]string `json:"errors"`
}

func (e APIErrors) Error() string {
	for _, m := range e.Messages {
		fmt.Println("[INFO] aaaaaaaaaaa", m)
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
