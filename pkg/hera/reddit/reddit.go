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
		log.Err(err).Msg("Error initializing reddit client")
		return Reddit{}, err
	}
	r.ReadOnly = client
	credentials := reddit.Credentials{ID: id, Secret: secret, Username: u, Password: pw}
	client, err = reddit.NewClient(credentials)
	if err != nil {
		log.Err(err).Msg("Error initializing reddit client")
		return Reddit{}, err
	}
	r.FullClient = client
	return r, err
}

func (r *Reddit) GetTopPosts(ctx context.Context, subreddit string, lpo *reddit.ListPostOptions) ([]*reddit.Post, *reddit.Response, error) {
	posts, resp, err := r.ReadOnly.Subreddit.TopPosts(ctx, subreddit, lpo)
	if err != nil {
		log.Err(err).Msg("Error getting top posts")
		return nil, nil, err
	}
	return posts, resp, nil
}

func (r *Reddit) GetControversialPosts(ctx context.Context, subreddit string, lpo *reddit.ListPostOptions) ([]*reddit.Post, *reddit.Response, error) {
	posts, resp, err := r.ReadOnly.Subreddit.ControversialPosts(ctx, subreddit, lpo)
	if err != nil {
		log.Err(err).Msg("Error getting top posts")
		return nil, nil, err
	}
	return posts, resp, nil
}

func (r *Reddit) GetNewPosts(ctx context.Context, subreddit string, lpo *reddit.ListOptions) ([]*reddit.Post, *reddit.Response, error) {
	posts, resp, err := r.ReadOnly.Subreddit.NewPosts(ctx, subreddit, lpo)
	if err != nil {
		log.Err(err).Msg("Error getting new posts")
		return nil, nil, err
	}
	return posts, resp, nil
}

func (r *Reddit) GetRisingPosts(ctx context.Context, subreddit string, lpo *reddit.ListOptions) ([]*reddit.Post, *reddit.Response, error) {
	posts, resp, err := r.ReadOnly.Subreddit.RisingPosts(ctx, subreddit, lpo)
	if err != nil {
		log.Err(err).Msg("Error getting rising posts")
		return nil, nil, err
	}
	return posts, resp, nil
}

func (r *Reddit) GetSubreddit(ctx context.Context, subreddit string) (*reddit.Subreddit, *reddit.Response, error) {
	sr, resp, err := r.ReadOnly.Subreddit.Get(ctx, subreddit)
	if err != nil {
		log.Err(err).Msg("Error getting top posts")
		return nil, nil, err
	}
	return sr, resp, nil
}
