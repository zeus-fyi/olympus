package hera_v1_codegen

import (
	"context"
	"fmt"
	"strings"

	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

const (
	DbSchemaDir     = "datastores/postgres/local_docker/docker-entrypoint-initdb.d/init.sql"
	HeraDbModelsDir = "datastores/postgres/apps/hera/models"
	PkgDir          = "pkg"
	AppsDir         = "apps"
	CookbooksDir    = "cookbooks"
	DockerDir       = "docker"
	DatastoresDir   = "datastores"
	AiZeusApi       = "apps/olympus/zeus/api/v1/zeus/ai"
)

type BuildAiInstructions struct {
	Path                filepaths.Path
	PromptInstructions  string
	OrderedInstructions []BuildAiFileInstruction
	FileReferencesMap   map[string]CodeFilesMetadata
	SearchPath          map[string][]string
}

type BuildAiFileInstruction struct {
	DirIn                           string
	FileName                        string
	FileLevelInstruction            string
	OrderedFileFunctionInstructions []FunctionInstruction
	OrderedGoTypeInstructions       []GoTypeInstruction
}

type GoTypeInstruction struct {
	GoTypeInstruction string
	GoTypeName        string
	GoType            string
}

func (b *BuildAiInstructions) SetSearchPath() {
	m := make(map[string][]string)
	for i, v := range b.OrderedInstructions {
		if i == 0 {
			m[v.DirIn] = []string{}
		}
		m[v.DirIn] = append(m[v.DirIn], v.FileName)
	}
	b.SearchPath = m
}

type FunctionInstruction struct {
	FunctionInstruction string
	FunctionInfo        FunctionInfo
}

func FormatInstructionPrompt(instructions string) string {
	return fmt.Sprintf("START PROMPT: %s\nEND PROMPT\n", instructions)
}

func GenerateInstructions(ctx context.Context, bai *BuildAiInstructions) string {
	if bai == nil {
		return ""
	}
	bai.SetSearchPath()
	scMap, err := ExtractSourceCode(ctx, bai)
	if err != nil {
		panic(err)
	}

	if scMap == nil {
		return ""
	}
	prompt := bai.PromptInstructions
	for _, v := range bai.OrderedInstructions {
		fr, ok := bai.FileReferencesMap[v.DirIn]
		if !ok {
			continue
		}

		sqlFile, ok := fr.SQLCodeFiles.Files[v.FileName]
		if ok {
			prompt += fmt.Sprintf("%s \nFilepath: %s/%s \n%s\n", v.FileLevelInstruction, v.DirIn, v.FileName, sqlFile.Contents)
			continue
		}

		goFile, ok := fr.GoCodeFiles.Files[v.FileName]
		if !ok {
			continue
		}
		prompt += fmt.Sprintf("%s \nFilepath: %s/%s \n", v.FileLevelInstruction, v.DirIn, v.FileName)
		prompt += fmt.Sprintf("Imports: %s\n", goFile.Imports)
		for _, f := range v.OrderedFileFunctionInstructions {
			funcInfo, aok := goFile.Functions[f.FunctionInfo.Name]
			if !aok {
				continue
			}
			switch f.FunctionInstruction {
			default:
				prompt += f.FunctionInfo.Name + ": " + f.FunctionInstruction + "\n" + " func parameters: " + funcInfo.Parameters + " func return type:" + funcInfo.ReturnType + " body: " + funcInfo.Body + "\n"
			}
		}
		for _, gt := range v.OrderedGoTypeInstructions {
			switch strings.ToLower(gt.GoType) {
			case "struct":
				gtInfo, aok := goFile.Structs[gt.GoType]
				if !aok {
					continue
				}
				prompt += gt.GoTypeInstruction + "\n" + gtInfo + "\n"
			}

		}
	}
	return FormatInstructionPrompt(prompt)
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
