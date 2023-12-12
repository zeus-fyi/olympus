package hera_reddit

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type Reddit struct {
	ReadOnly   *reddit.Client
	FullClient *reddit.Client
}

var RedditClient Reddit

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
	RedditClient.ReadOnly = r.ReadOnly
	RedditClient.FullClient = client
	r.FullClient = client
	return r, err
}

func InitOrgRedditClient(ctx context.Context, id, secret, u, pw string) (Reddit, error) {
	r := Reddit{}
	userAgent := fmt.Sprintf("Zeusfyi/%s (by /u/%s)", id, u)
	log.Info().Str("userAgent", userAgent)
	ro, err := reddit.NewReadonlyClient(reddit.WithUserAgent(userAgent))
	if err != nil {
		log.Err(err).Interface("userAgent", userAgent).Msg("Error initializing reddit client")
		return Reddit{}, err
	}
	r.ReadOnly = ro
	credentials := reddit.Credentials{ID: id, Secret: secret, Username: u, Password: pw}
	fc, err := reddit.NewClient(credentials)
	if err != nil {
		log.Err(err).Interface("userAgent", userAgent).Msg("Error initializing reddit client")
		return Reddit{}, err
	}
	r.FullClient = fc
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

type RedditPostSearchResponse struct {
	Posts []*reddit.Post   `json:"posts"`
	Resp  *reddit.Response `json:"resp"`
}

func (r *Reddit) GetNewPosts(ctx context.Context, subreddit string, lpo *reddit.ListOptions) (*RedditPostSearchResponse, error) {
	posts, resp, err := r.ReadOnly.Subreddit.NewPosts(ctx, subreddit, lpo)
	if err != nil {
		log.Err(err).Interface("resp", resp).Msg("Error getting new posts")
		return nil, err
	}
	re := &RedditPostSearchResponse{
		Posts: posts,
		Resp:  resp,
	}
	return re, nil
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
