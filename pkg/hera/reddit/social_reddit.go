package hera_reddit

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type SocialReddit struct {
	RedditReplyToPost RedditReplyToPost `json:"redditReplyToPost"`
}
type RedditReplyToPost struct {
	ParentID string `json:"parentID"`
	Reply    string `json:"reply"`
}

func (r *Reddit) ReplyToPost(ctx context.Context, parentID, reply string) (*reddit.Comment, *reddit.Response, error) {
	cm, cr, err := r.FullClient.Comment.Submit(ctx, parentID, reply)
	if err != nil {
		log.Err(err).Msg("Error submitting comment")
		return nil, nil, err
	}
	return cm, cr, nil
}
