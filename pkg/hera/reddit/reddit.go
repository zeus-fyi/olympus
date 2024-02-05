package hera_reddit

import (
	"context"
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/vartanbeno/go-reddit/v2/reddit"
	resty_base "github.com/zeus-fyi/zeus/zeus/z_client/base"
)

type Reddit struct {
	ReadOnly   *reddit.Client
	FullClient *reddit.Client

	Resty       resty_base.Resty
	AccessToken string
}

var RedditClient Reddit

func CreateFormattedStringRedditUA(platform, appID, versionString, redditUsername string) string {
	return fmt.Sprintf("%s:%s:%s (by /u/%s)", platform, appID, versionString, redditUsername)
}

func InitRedditClient(ctx context.Context, id, secret, u, pw string) (Reddit, error) {
	r := Reddit{}
	ro, err := reddit.NewReadonlyClient(reddit.WithUserAgent(CreateFormattedStringRedditUA("web", "zeusfyi", "0.0.1", "zeus-fyi")))
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
	ro, err := reddit.NewReadonlyClient(reddit.WithUserAgent(CreateFormattedStringRedditUA("web", "zeusfyi", "0.0.1", "zeus-fyi")))
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
	r.AccessToken = respToken.AccessToken
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
		Children []*Child `json:"children"`
	} `json:"data"`
}

type Child struct {
	Data *reddit.Post `json:"data"`
}

func (r *Reddit) GetRedditReq(ctx context.Context, path string, responseMap *echo.Map, urlVals url.Values) (*resty.Response, error) {
	ua := CreateFormattedStringRedditUA("web", "zeusfyi", "0.0.1", "zeus-fyi")
	r.Resty.SetHeader("User-Agent", ua)
	r.Resty.QueryParam = urlVals
	resp, err := r.Resty.R().SetResult(responseMap).Get(path)
	if err != nil {
		log.Err(err).Interface("resp", resp).Msg("Error getting new posts")
		return nil, err
	}
	if resp.StatusCode() >= 400 {
		log.Err(err).Interface("resp", resp).Msg("Error getting new posts")
		return nil, fmt.Errorf("error getting new posts")
	}
	return resp, nil
}

func (r *Reddit) PostRedditReq(ctx context.Context, path string, payload echo.Map, responseMap *echo.Map) (*resty.Response, error) {
	ua := CreateFormattedStringRedditUA("web", "zeusfyi", "0.0.1", "zeus-fyi")
	r.Resty.SetHeader("User-Agent", ua)
	resp, err := r.Resty.R().SetBody(&payload).SetResult(responseMap).Post(path)
	if err != nil {
		log.Err(err).Interface("resp", resp).Msg("Error getting new posts")
		return nil, err
	}
	if resp.StatusCode() >= 400 {
		log.Err(err).Interface("resp", resp).Msg("Error getting new posts")
		return nil, fmt.Errorf("error getting new posts")
	}
	return resp, nil
}

func (r *Reddit) PutRedditReq(ctx context.Context, path string, payload echo.Map, responseMap *echo.Map) (*resty.Response, error) {
	ua := CreateFormattedStringRedditUA("web", "zeusfyi", "0.0.1", "zeus-fyi")
	r.Resty.SetHeader("User-Agent", ua)
	resp, err := r.Resty.R().SetBody(&payload).SetResult(responseMap).Put(path)
	if err != nil {
		log.Err(err).Interface("resp", resp).Msg("Error getting new posts")
		return nil, err
	}
	if resp.StatusCode() >= 400 {
		log.Err(err).Interface("resp", resp).Msg("Error getting new posts")
		return nil, fmt.Errorf("error getting new posts")
	}
	return resp, nil
}

func (r *Reddit) GetMe(ctx context.Context) (*reddit.User, error) {
	path := fmt.Sprintf("/api/v1/me")
	ua := CreateFormattedStringRedditUA("web", "zeusfyi", "0.0.1", "zeus-fyi")
	r.Resty.SetHeader("User-Agent", ua)
	user := &reddit.User{}
	resp, err := r.Resty.R().SetResult(&user).Get(path)
	if err != nil {
		log.Err(err).Interface("resp", resp).Msg("Error getting new posts")
		return nil, err
	}
	if resp.StatusCode() >= 400 {
		log.Err(err).Interface("resp", resp).Msg("Error getting new posts")
		return nil, fmt.Errorf("error getting new posts")
	}
	return user, nil
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
	if lpo.After == "" {
		path = fmt.Sprintf("/r/%s/new.json?limit=100", subreddit)
	}
	ua := CreateFormattedStringRedditUA("web", "zeusfyi", "0.0.1", "zeus-fyi")
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
	var posts []*reddit.Post
	for _, c := range s.Data.Children {
		posts = append(posts, c.Data)
	}
	return posts, nil
}

func (r *Reddit) GetLastLikedPostV2(ctx context.Context, username string) ([]*reddit.Post, error) {
	path := fmt.Sprintf("user/%s/upvoted?limit=1", username)
	ua := CreateFormattedStringRedditUA("web", "zeusfyi", "0.0.1", "zeus-fyi")
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
	var posts []*reddit.Post
	for _, c := range s.Data.Children {
		posts = append(posts, c.Data)
	}
	return posts, nil
}

func (r *Reddit) GetLastLikedPost(ctx context.Context, userName string) ([]*reddit.Post, error) {
	// Assuming you have a method to set up OAuth2 and headers
	//path := fmt.Sprintf("/user/%s/liked?limit=1", userName) // Replace {username} with your actual username
	opt := &reddit.ListUserOverviewOptions{
		ListOptions: reddit.ListOptions{
			Limit: 1,
		},
	}
	resp, _, err := r.FullClient.User.UpvotedOf(ctx, userName, opt)
	if err != nil {
		log.Err(err).Interface("resp", resp).Msg("Error getting liked posts")
		return nil, err
	}
	return resp, nil
}

func (r *Reddit) PostComment(ctx context.Context, parentID, textReply string) (*reddit.Comment, error) {
	resp, _, err := r.FullClient.Comment.Submit(ctx, parentID, textReply)
	if err != nil {
		log.Err(err).Interface("resp", resp).Msg("Error getting liked posts")
		return nil, err
	}
	return resp, nil
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
