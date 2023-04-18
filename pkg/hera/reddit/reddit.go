package hera_reddit

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type Reddit struct {
	ReadOnly   *reddit.Client
	FullClient *reddit.Client
}

func InitRedditClient(ctx context.Context, id, secret, u, pw string) (Reddit, error) {
	r := Reddit{}
	client, err := reddit.NewReadonlyClient(reddit.WithUserAgent("Zeusfyi/1.0 (by /u/zeus-fyi)"))
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("Error initializing reddit client")
		return Reddit{}, err
	}
	r.ReadOnly = client
	credentials := reddit.Credentials{ID: id, Secret: secret, Username: u, Password: pw}
	client, err = reddit.NewClient(credentials)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("Error initializing reddit client")
		return Reddit{}, err
	}
	r.FullClient = client
	return r, err
}
