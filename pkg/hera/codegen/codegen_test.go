package hera_v1_codegen

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/test"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

var (
	dirIn = "../../.."
)

type CodeGenTestSuite struct {
	test.AutoGenBaseTestSuiteBase
}

func (s *CodeGenTestSuite) SetupTest() {
	s.InitLocalConfigs()
	UseAutoGenDirectory()
}

var (
	ctx = context.Background()
)

func (s *CodeGenTestSuite) TestCreateAiAssistantCodeGenWorkflowInstructions2() {
	f := filepaths.Path{
		DirIn:       dirIn,
		FilterFiles: sf,
	}
	actInst := ``
	bins := BuildAiInstructions{
		Path: f,
		OrderedInstructions: []BuildAiFileInstruction{
			{
				DirIn:                PkgDir + "/zeus/ai/orchestrations",
				FileName:             "workflows_twitter.go",
				FileLevelInstruction: actInst,
				OrderedFileFunctionInstructions: []FunctionInstruction{
					{
						FunctionInstruction: "Use this function as an example reference for adding a new activity to the workflow",
						FunctionInfo: FunctionInfo{
							Name: "AiIngestTwitterWorkflow",
						},
					},
				},
			},
			{
				DirIn:                PkgDir + "/zeus/ai/orchestrations",
				FileName:             "activities.go",
				FileLevelInstruction: actInst,
				OrderedFileFunctionInstructions: []FunctionInstruction{
					{
						FunctionInstruction: "Use this activity function for adding a new activity to the workflow",
						FunctionInfo: FunctionInfo{
							Name: "SearchRedditNewPostsUsingSubreddit",
						},
					},
				},
			},
			{
				DirIn:                PkgDir + "/zeus/ai/orchestrations",
				FileName:             "workflows_reddit.go",
				FileLevelInstruction: "",
				OrderedFileFunctionInstructions: []FunctionInstruction{
					{
						FunctionInstruction: "Add the activity to the workflow section after the UpsertAssignmentActivity",
						FunctionInfo: FunctionInfo{
							Name: "AiIngestRedditWorkflow",
						},
					},
				},
			},
			{
				DirIn:    PkgDir + "/hera/reddit",
				FileName: "reddit.go",
				OrderedGoTypeInstructions: []GoTypeInstruction{
					{
						GoTypeInstruction: "create a var with this struct type that gets assigned from the output of the activity",
						GoType:            "struct",
						GoTypeName:        "RedditPostSearchResponse",
					},
				},
			},
		},
	}
	prompt := GenerateInstructions(ctx, &bins)
	fmt.Println(prompt)

	hera_openai.InitHeraOpenAI(s.Tc.OpenAIAuth)
	params := hera_openai.OpenAIParams{
		Prompt: prompt,
	}
	ou := org_users.NewOrgUserWithID(s.Tc.ProductionLocalTemporalOrgID, s.Tc.ProductionLocalTemporalUserID)
	resp, err := hera_openai.HeraOpenAI.MakeCodeGenRequestV2(ctx, ou, params)
	s.Require().NoError(err)
	fmt.Println(resp.Choices[0].Message.Content)
	f.DirOut = "./generated_outputs"
	f.FnOut = "workflow_instructions.txt"
	err = f.WriteToFileOutPath([]byte(prompt))
	s.Require().NoError(err)
}

func (s *CodeGenTestSuite) TestCodeGenFunction() {
	f := filepaths.Path{
		DirIn:       dirIn,
		DirOut:      "./",
		FnOut:       "codegen_output.json",
		FilterFiles: sf,
	}
	prompt := GenerateSqlTableFromExample(f)
	s.Require().NotEmpty(prompt)
	fmt.Println(prompt)

	hera_openai.InitHeraOpenAI(s.Tc.OpenAIAuth)
	params := hera_openai.OpenAIParams{
		Prompt: prompt,
	}
	ou := org_users.NewOrgUserWithID(s.Tc.ProductionLocalTemporalOrgID, s.Tc.ProductionLocalTemporalUserID)
	resp, err := hera_openai.HeraOpenAI.MakeCodeGenRequestV2(ctx, ou, params)
	s.Require().NoError(err)
	fmt.Println(resp.Choices[0].Message.Content)

	s.Require().NoError(err)

	err = f.WriteToFileOutPath([]byte(prompt))
	s.Require().NoError(err)

}

func (s *CodeGenTestSuite) TestCreateAiAssistantCodeGenWorkflowInstructions() {
	f := filepaths.Path{
		DirIn:       dirIn,
		FilterFiles: sf,
	}
	actInst := `write: add a new activity name that matches the existing syntax style for Searching New SubReddit Posts,
					the name should be derived from the reference func name in reddit.go
					and then create a new func that matches the new activity name and
				write the logic for the new wrapper func, which calls the reference func from reddit.go and uses RedditClient from reddit.go to make the call`
	bins := BuildAiInstructions{
		Path: f,
		OrderedInstructions: []BuildAiFileInstruction{
			{
				DirIn:                PkgDir + "/zeus/ai/orchestrations",
				FileName:             "activities.go",
				FileLevelInstruction: actInst,
				OrderedFileFunctionInstructions: []FunctionInstruction{
					{
						FunctionInstruction: "Add the new activity definition here to this function, this is an struct pointer function, and you need to add h. as the prefix to any new activity definition",
						FunctionInfo: FunctionInfo{
							Name: "GetActivities",
						},
					},
					{
						FunctionInstruction: "Read Only Reference syntax style for Building Searching New SubReddit Post Activity",
						FunctionInfo: FunctionInfo{
							Name: "SearchTwitterUsingQuery",
						},
					},
				},
			},
			{
				DirIn:                PkgDir + "/hera/reddit",
				FileName:             "reddit.go",
				FileLevelInstruction: "",
				OrderedFileFunctionInstructions: []FunctionInstruction{{
					FunctionInstruction: "Read Only Reference syntax style for Building Searching New SubReddit Post Activity",
					FunctionInfo: FunctionInfo{
						Name: "GetNewPosts",
					},
				}},
			},
		},
	}
	//for _, is := range bins.OrderedInstructions {
	//	fmt.Println(is.FileLevelInstruction)
	//}

	prompt := GenerateInstructions(ctx, &bins)
	fmt.Println(prompt)

	hera_openai.InitHeraOpenAI(s.Tc.OpenAIAuth)
	params := hera_openai.OpenAIParams{
		Prompt: prompt,
	}
	ou := org_users.NewOrgUserWithID(s.Tc.ProductionLocalTemporalOrgID, s.Tc.ProductionLocalTemporalUserID)
	resp, err := hera_openai.HeraOpenAI.MakeCodeGenRequestV2(ctx, ou, params)
	s.Require().NoError(err)
	fmt.Println(resp.Choices[0].Message.Content)
}

func (s *CodeGenTestSuite) TestCreateCodeSourceParsing() {
	f := filepaths.Path{
		DirIn:       dirIn,
		FilterFiles: sf,
	}
	bins := &BuildAiInstructions{
		Path: f,
	}
	b, err := ExtractSourceCode(ctx, bins)
	s.NoError(err)
	s.NotEmpty(b)

	//tmp := b.FileReferencesMap[DbSchemaDir]
	//for _, fvs := range tmp.SQLCodeFiles.Files {
	//	fmt.Println(fvs.FileName)
	//}

	//directoryPath := PkgDir + "/zeus/ai/orchestrations"
	//fmt.Println("Directory Path: ", directoryPath)
	//goCode := b.FileReferencesMap[PkgDir+"/zeus/ai/orchestrations"]
	//for _, fvs := range goCode.GoCodeFiles.Files {
	//	fmt.Println(fvs.FileName)
	//}
	//
	//fmt.Println("Directory Imports...")
	//for _, di := range goCode.GoCodeFiles.DirectoryImports {
	//	fmt.Println(di)
	//}
	jsCode := b.FileReferencesMap["apps/olympus/hestia/assets/src/app"]
	fmt.Println("Js/Tsx Imports...")
	for _, di := range jsCode.JsCodeFiles.Files {
		fmt.Println(di.FileName)
	}
}

func TestCodeGenTestSuite(t *testing.T) {
	suite.Run(t, new(CodeGenTestSuite))
}
func UseAutoGenDirectory() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "")
	err := os.Chdir(dir)
	if err != nil {
		panic(err.Error())
	}
	return dir
}

/*

Generated output

```go
// pkg/zeus/ai/orchestrations/activities.go

func GetActivities() []interface{} {
	ka := kronos_helix.NewKronosActivities()
	actSlice := []interface{}{
		h.AiTask, h.SaveAiTaskResponse, h.SendTaskResponseEmail, h.InsertEmailIfNew,
		h.InsertAiResponse, h.InsertTelegramMessageIfNew,
		h.InsertIncomingTweetsFromSearch, h.SearchTwitterUsingQuery, h.SelectTwitterSearchQuery,
		h.SearchNewSubRedditPosts,
	}

	return append(actSlice, ka.GetActivities()...)
}

func SearchNewSubRedditPosts(ctx context.Context, subreddit string, lpo RedditListPostOptions) ([]*reddit.Post, *reddit.Response, error) {
	posts, resp, err := hera_reddit.RedditClient.GetNewPosts(ctx, subreddit, lpo)
	if err != nil {
		log.Err(err).Msg("SearchNewSubRedditPosts")
		return nil, nil, err
	}
	return posts, resp, nil
}

// pkg/hera/reddit/reddit.go

// Add this struct to the file if it's not already defined
type RedditListPostOptions struct {
	ListOptions reddit.ListOptions
	Time        string
}

var RedditClient *reddit.Client

// Add the RedditClient initialization somewhere in your codebase, if not already present.
*/
