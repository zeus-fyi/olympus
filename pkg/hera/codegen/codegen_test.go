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
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
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

func (s *CodeGenTestSuite) TestCreateWorkflow() {
	ctx := context.Background()
	sf := &strings_filter.FilterOpts{
		DoesNotStartWithThese: []string{"configs", "sandbox", "apps/external", ".git", ".circleci", ".DS_Store", ".idea", "apps/zeus/test/configs", "pkg/.DS_Store"},
		//StartsWithAnyOfThese:  []string{"apps", "pkg", "docker", ".github", "cookbooks", "datastores", "datastores/postgres/local_docker/docker-entrypoint-initdb.d/init.sql"},
		DoesNotInclude: []string{"hardhat/artifacts", "node_modules", ".kube", "bin", "build", ".git", "hardhat/cache"},
	}
	f := filepaths.Path{
		DirIn:       dirIn,
		FilterFiles: sf,
	}
	b, err := CreateWorkflow(ctx, f)
	s.NoError(err)
	s.NotEmpty(b)

	//tmp := b.Map[DbSchemaDir]
	//for _, fvs := range tmp.SQLCodeFiles.Files {
	//	fmt.Println(fvs.FileName)
	//}

	directoryPath := PkgDir + "/zeus/ai/orchestrations"
	fmt.Println("Directory Path: ", directoryPath)
	goCode := b.Map[PkgDir+"/zeus/ai/orchestrations"]
	for _, fvs := range goCode.GoCodeFiles.Files {
		fmt.Println(fvs.FileName)
	}

	fmt.Println("Directory Imports...")
	for _, di := range goCode.GoCodeFiles.DirectoryImports {
		fmt.Println(di)
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
