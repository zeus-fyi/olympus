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
		StartsWithAnyOfThese:  []string{"apps", "pkg", "docker", ".github", "cookbooks", "datastores"},
		DoesNotInclude:        []string{"hardhat/artifacts", "node_modules", ".kube", "bin", "build", ".git", "hardhat/cache"},
	}
	sf.DoesNotInclude = append(sf.DoesNotInclude, []string{"go-ethereum", "apps/external/tables-to-go", "tmp", "vendor", "td", "tojen"}...)
	f := filepaths.Path{
		DirIn:       dirIn,
		FilterFiles: sf,
	}
	b, err := CreateWorkflow(ctx, f)
	s.NoError(err)
	s.NotEmpty(b)
	fmt.Println(string(b))
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
