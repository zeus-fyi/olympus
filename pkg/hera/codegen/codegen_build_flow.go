package hera_v1_codegen

import (
	"context"
	"fmt"

	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

const (
	DbSchemaDir   = "datastores/postgres/local_docker/docker-entrypoint-initdb.d/init.sql"
	PkgDir        = "pkg"
	AppsDir       = "apps"
	CookbooksDir  = "cookbooks"
	DockerDir     = "docker"
	DatastoresDir = "datastores"
	AiZeusApi     = "apps/olympus/zeus/api/v1/zeus/ai"
)

type BuildAiInstructions struct {
	Instructions []BuildAiInstruction
}

type BuildAiInstruction struct {
	DirIn               string
	FileInstructionsMap map[string]string
}

type FileInfo struct {
	Functions []FunctionInfo
}

type GeneratedBuildAiInstructions struct {
	Instructions []GeneratedBuildAiInstruction
}
type GeneratedBuildAiInstruction struct {
	DirIn               string
	FileInstructionsMap map[string]string
}

func FormatInstructionPrompt(instructions, fileContents string) string {
	return fmt.Sprintf("START PROMPT: %s\n END PROMPT\n %s \n", instructions, fileContents)
}

func BuildAiInstructionsFromSourceCode(ctx context.Context, f filepaths.Path, buildDirs BuildAiInstructions) *GeneratedBuildAiInstructions {
	sc, err := ExtractSourceCode(ctx, f)
	if err != nil {
		panic(err)
	}
	gbi := &GeneratedBuildAiInstructions{}
	for _, bd := range buildDirs.Instructions {
		if bd.FileInstructionsMap == nil {
			continue
		}
		fmt.Println("dirIn: ", bd.DirIn)
		for _, gf := range sc.Map[bd.DirIn].GoCodeFiles.Files {
			if bd.FileInstructionsMap[gf.FileName] == "" {
				delete(bd.FileInstructionsMap, gf.FileName)
			} else {
				gbi.Instructions = append(gbi.Instructions, GeneratedBuildAiInstruction{
					DirIn:               bd.DirIn,
					FileInstructionsMap: map[string]string{gf.FileName: FormatInstructionPrompt(bd.FileInstructionsMap[gf.FileName], gf.Contents)},
				})
			}
		}
	}
	return gbi
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

// Subreddit holds information about a subreddit
type Subreddit struct {
	ID      string     `json:"id,omitempty"`
	FullID  string     `json:"name,omitempty"`
	Created *Timestamp `json:"created_utc,omitempty"`

	URL                  string `json:"url,omitempty"`
	Name                 string `json:"display_name,omitempty"`
	NamePrefixed         string `json:"display_name_prefixed,omitempty"`
	Title                string `json:"title,omitempty"`
	Description          string `json:"public_description,omitempty"`
	Type                 string `json:"subreddit_type,omitempty"`
	SuggestedCommentSort string `json:"suggested_comment_sort,omitempty"`

	Subscribers     int  `json:"subscribers"`
	ActiveUserCount *int `json:"active_user_count,omitempty"`
	NSFW            bool `json:"over18"`
	UserIsMod       bool `json:"user_is_moderator"`
	Subscribed      bool `json:"user_is_subscriber"`
	Favorite        bool `json:"user_has_favorited"`
}
*/
