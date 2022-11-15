package zeus

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/hera/lib"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/test"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

type ZeusRecipeTestSuite struct {
	test.AutoGenBaseTestSuiteBase
}

func createTestCodeGenShell() lib.CodeGen {
	p := structs.Path{
		PackageName: "autogen_structs",
		DirIn:       "",
		DirOut:      "tmp",
		FnIn:        "zeus.go",
		Env:         "",
	}
	cg := lib.NewCodeGen(p)
	return cg
}

func (s *ZeusRecipeTestSuite) TestZeusDerivativeStructGen() {
}

func TestZeusRecipeTestSuite(t *testing.T) {
	suite.Run(t, new(ZeusRecipeTestSuite))
}
