package gettr

import (
	"github.com/dghubble/sling"
	"net/http"
)

const gettrAPI = "https://api.gettr.com/"

// Client is a Gettr Client
type Client struct {
	sling *sling.Sling
	User  *UserService
}

// NewClient returns a new client
func NewClient(httpClient *http.Client) *Client {
	base := sling.New().Client(httpClient).Base(gettrAPI)
	return &Client{
		sling: base,
		User:  newUserService(base.New()),
	}
}
