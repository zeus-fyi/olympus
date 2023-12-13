package hera_reddit

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/vartanbeno/go-reddit/v2/reddit"
	resty_base "github.com/zeus-fyi/zeus/zeus/z_client/base"
)

type Reddit struct {
	ReadOnly   *reddit.Client
	FullClient *reddit.Client

	Resty resty_base.Resty
}

var RedditClient Reddit

func createFormattedString(platform, appID, versionString, redditUsername string) string {
	return fmt.Sprintf("%s:%s:%s (by /u/%s)", platform, appID, versionString, redditUsername)
}

func InitRedditClient(ctx context.Context, id, secret, u, pw string) (Reddit, error) {
	r := Reddit{}
	ro, err := reddit.NewReadonlyClient(reddit.WithUserAgent(createFormattedString("web", "zeusfyi", "0.0.1", "zeus-fyi")))
	if err != nil {
		log.Err(err).Msg("Error initializing reddit client")
		return Reddit{}, err
	}
	r.ReadOnly = ro
	credentials := reddit.Credentials{ID: id, Secret: secret, Username: u, Password: pw}
	client, err := reddit.NewClient(credentials)
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
	ro, err := reddit.NewReadonlyClient(reddit.WithUserAgent(createFormattedString("web", "zeusfyi", "0.0.1", "zeus-fyi")))
	if err != nil {
		log.Err(err).Msg("Error initializing reddit client")
		return Reddit{}, err
	}
	r.ReadOnly = ro
	credentials := reddit.Credentials{ID: id, Secret: secret, Username: u, Password: pw}
	fc, err := reddit.NewClient(credentials)
	if err != nil {
		log.Err(err).Msg("Error initializing reddit client")
		return Reddit{}, err
	}
	r.FullClient = fc

	rb := resty.New()
	payload := fmt.Sprintf("grant_type=password&username=%s&password=%s", u, pw)
	respToken := RedditOauth2Resp{}
	// Make a POST request
	resp, err := rb.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBasicAuth(id, secret). // Use 'secret', not 'pw'
		SetBody(payload).
		SetResult(&respToken).
		Post("https://www.reddit.com/api/v1/access_token")

	if err != nil {
		log.Err(err).Interface("resp", resp).Msg("Error getting access token")
		return Reddit{}, err
	}

	if resp.StatusCode() >= 400 {
		panic("Error getting access token")
	}
	rb.SetBaseURL("https://oauth.reddit.com")
	rb.SetAuthToken(respToken.AccessToken)
	r.Resty = resty_base.GetBaseRestyClient("https://oauth.reddit.com", respToken.AccessToken)
	return r, err
}

type RedditOauth2Resp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}
type RedditResponse struct {
	Data struct {
		Children []*reddit.Post `json:"children"`
	} `json:"data"`
}

func (r *Reddit) GetNewPosts(ctx context.Context, subreddit string, lpo *reddit.ListOptions) ([]*reddit.Post, error) {
	if lpo == nil {
		lpo = &reddit.ListOptions{
			Limit:  100,
			After:  "",
			Before: "",
		}
	}
	path := fmt.Sprintf("/r/%s/new.json?limit=100&after=%s", subreddit, lpo.After)
	ua := createFormattedString("web", "zeusfyi", "0.0.1", "zeus-fyi")
	r.Resty.SetHeader("User-Agent", ua)
	var s RedditResponse
	resp, err := r.Resty.R().SetResult(&s).Get(path)
	if err != nil {
		log.Err(err).Interface("resp", resp).Msg("Error getting new posts")
		return nil, err
	}
	if resp.StatusCode() >= 400 {
		log.Err(err).Interface("resp", resp).Msg("Error getting new posts")
		return nil, fmt.Errorf("error getting new posts")
	}
	return s.Data.Children, nil
}

func (r *Reddit) GetTopPosts(ctx context.Context, subreddit string, lpo *reddit.ListPostOptions) ([]*reddit.Post, *reddit.Response, error) {
	posts, resp, err := r.FullClient.Subreddit.TopPosts(ctx, subreddit, lpo)
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
