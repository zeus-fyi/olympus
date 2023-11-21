package hera_v1_codegen

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/test"
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
	actInst := `update func (h *ZeusAiPlatformActivities) GetActivities() ActivitiesSlice
				write: add a new activity name that matches the existing syntax style for Searching New SubReddit Posts,
					the name should be derived from the reference func name in reddit.go
					and then create a new func that matches the new activity name and
				write the logic for the new wrapper func, which calls the reference func from reddit.go and uses RedditClient from reddit.go to make the call`

	bi := BuildAiInstructions{
		Instructions: []BuildAiInstruction{
			{
				DirIn: PkgDir + "/hera/reddit",
				FileInstructionsMap: map[string]string{
					"reddit.go": `reference func (r *Reddit) GetNewPosts(ctx context.Context, subreddit string, lpo *reddit.ListOptions)
							      reference var RedditClient Reddit`,
				},
			},
			{
				DirIn: PkgDir + "/zeus/ai/orchestrations",
				FileInstructionsMap: map[string]string{
					"activities.go": actInst,
				},
			},
		},
	}

	gbi := BuildAiInstructionsFromSourceCode(ctx, f, bi)
	for _, g := range gbi.Instructions {
		fmt.Println(g.DirIn)
		for k, v := range g.FileInstructionsMap {
			fmt.Println(k, " | ", v)
		}
	}
}

func (s *CodeGenTestSuite) TestCreateCodeSourceParsing() {
	f := filepaths.Path{
		DirIn:       dirIn,
		FilterFiles: sf,
	}
	b, err := ExtractSourceCode(ctx, f)
	s.NoError(err)
	s.NotEmpty(b)

	//tmp := b.Map[DbSchemaDir]
	//for _, fvs := range tmp.SQLCodeFiles.Files {
	//	fmt.Println(fvs.FileName)
	//}

	//directoryPath := PkgDir + "/zeus/ai/orchestrations"
	//fmt.Println("Directory Path: ", directoryPath)
	//goCode := b.Map[PkgDir+"/zeus/ai/orchestrations"]
	//for _, fvs := range goCode.GoCodeFiles.Files {
	//	fmt.Println(fvs.FileName)
	//}
	//
	//fmt.Println("Directory Imports...")
	//for _, di := range goCode.GoCodeFiles.DirectoryImports {
	//	fmt.Println(di)
	//}
	jsCode := b.Map["apps/olympus/hestia/assets/src/app"]
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
