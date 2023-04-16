package hera_twitter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/g8rswimmer/go-twitter/v2"
)

type Twitter struct {
	*twitter.Client
}

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}
func InitTwitterClient(ctx context.Context, bearer string) Twitter {
	client := &twitter.Client{
		Authorizer: authorize{Token: bearer},
		Client:     http.DefaultClient,
		Host:       "https://api.twitter.com",
	}
	return Twitter{client}
}
