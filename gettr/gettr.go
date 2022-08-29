package gettr

import (
	"encoding/json"
	"github.com/dghubble/sling"
	"net/http"
)

const gettrAPI = "https://api.gettr.com/"

// Client is a Gettr Client
type Client struct {
	httpClient *http.Client
	sling      *sling.Sling
	User       *UserService
	username   string
	userID     string
	token      string
	authHeader string
}

type authHeader struct {
	User  string `json:"user"`
	Token string `json:"token"`
}

// NewClient returns a new client
func NewClient(httpClient *http.Client) *Client {
	base := sling.New().Client(httpClient).Base(gettrAPI)
	client := Client{sling: base}
	client.httpClient = httpClient
	client.User = newUserService(base.New(), &client)
	return &client
}

// SetAuthToken sets the authentication username, userId and token for the actions that require permissions
func (c *Client) SetAuthToken(username, userID, token string) {
	c.username = username
	c.userID = userID
	c.token = token

	bytes, err := json.Marshal(authHeader{User: userID, Token: token})
	if err == nil {
		c.authHeader = string(bytes)
	}
}
