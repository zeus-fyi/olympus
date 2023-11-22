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
	UseAutoGenDirectory()
}

var (
	ctx = context.Background()
)

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
						FunctionInstruction: "Add the new activity definition here to this function",
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
	resp, err := hera_openai.HeraOpenAI.MakeCodeGenRequest(ctx, ou, params)
	s.NoError(err)
	fmt.Println(resp.Choices[0].Text)
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
