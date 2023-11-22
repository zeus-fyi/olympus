package hera_v1_codegen

import (
	"context"

	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

func GenerateSqlTableFromExample(f filepaths.Path) string {
	actInst := `write: two new SQL table definitions called ai_reddit_search_query and ai_reddit_incoming_posts. we index this data and search the date range using the primary key id which
				is a unix timestamp, and filter on text from post body, subreddit, and title for the search query derive a solution based on how we used indexing for tweets
				the final solution should allow us to insert the below data into the new tables and then query the data using the below data as an example
		        group your generated answers using the filepath of the file that contains most likely the reference for the answer 
				type Post struct {
					// generate columns for the below fields, use bigint for every int field
					ID                   string
					FullID               string
					Created              *Timestamp
					Edited               *Timestamp
					Permalink            string
					URL                  string
					Title                string
					Body                 string
					Likes                *bool
					Score                int
					UpvoteRatio          float32
					NumberOfComments     int
					
					// all elements below this line should be stored in a jsonb column
					SubredditName        string
					SubredditNamePrefixed string
					SubredditID          string
					SubredditSubscribers int
					Author               string
					AuthorID             string
					Spoiler              bool
					Locked               bool
					NSFW                 bool
					IsSelfPost           bool
					Saved                bool
					Stickied             bool
				}`
	bins := &BuildAiInstructions{
		Path:               f,
		PromptInstructions: actInst,
		OrderedInstructions: []BuildAiFileInstruction{
			{
				DirIn:                DbSchemaDir,
				FileName:             "603_ai_twitter.sql",
				FileLevelInstruction: "Use the example SQL table definitions and indexes in this file to create the new SQL table definitions",
			},
			{
				DirIn:                HeraDbModelsDir + "/search",
				FileName:             "twitter.go",
				FileLevelInstruction: "Use the example insert and select query functions to build an equivalent for the new reddit SQL definitions you create",
				OrderedFileFunctionInstructions: []FunctionInstruction{
					{
						FunctionInstruction: "Reference 1",
						FunctionInfo: FunctionInfo{
							Name: "insertTwitterSearchQuery",
						},
					},
					{
						FunctionInstruction: "Reference 1",
						FunctionInfo: FunctionInfo{
							Name: "InsertTwitterSearchQuery",
						},
					},
					{
						FunctionInstruction: "Reference 2",
						FunctionInfo: FunctionInfo{
							Name: "selectTwitterSearchQuery",
						},
					},
					{
						FunctionInstruction: "Reference 2",
						FunctionInfo: FunctionInfo{
							Name: "SelectTwitterSearchQuery",
						},
					},
					{
						FunctionInstruction: "Reference 2",
						FunctionInfo: FunctionInfo{
							Name: "insertIncomingTweets",
						},
					},
					{
						FunctionInstruction: "Reference 2",
						FunctionInfo: FunctionInfo{
							Name: "InsertIncomingTweets",
						},
					},
				},
			},
		},
	}
	prompt := GenerateInstructions(context.Background(), bins)
	return prompt
}

/*
// Post is a submitted post on Reddit.
type Post struct {
	ID      string     `json:"id,omitempty"`
	FullID  string     `json:"name,omitempty"`
	Created *Timestamp `json:"created_utc,omitempty"`
	Edited  *Timestamp `json:"edited,omitempty"`

	Permalink string `json:"permalink,omitempty"`
	URL       string `json:"url,omitempty"`

	Title string `json:"title,omitempty"`
	Body  string `json:"selftext,omitempty"`

	// Indicates if you've upvoted/downvoted (true/false).
	// If neither, it will be nil.
	Likes *bool `json:"likes"`

	Score            int     `json:"score"`
	UpvoteRatio      float32 `json:"upvote_ratio"`
	NumberOfComments int     `json:"num_comments"`

	SubredditName         string `json:"subreddit,omitempty"`
	SubredditNamePrefixed string `json:"subreddit_name_prefixed,omitempty"`
	SubredditID           string `json:"subreddit_id,omitempty"`
	SubredditSubscribers  int    `json:"subreddit_subscribers"`

	Author   string `json:"author,omitempty"`
	AuthorID string `json:"author_fullname,omitempty"`

	Spoiler    bool `json:"spoiler"`
	Locked     bool `json:"locked"`
	NSFW       bool `json:"over_18"`
	IsSelfPost bool `json:"is_self"`
	Saved      bool `json:"saved"`
	Stickied   bool `json:"stickied"`
}
*/
