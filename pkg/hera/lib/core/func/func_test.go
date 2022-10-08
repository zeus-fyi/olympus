package _func

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/core/primitives"
)

type FuncTestSuite struct {
	suite.Suite
}

func (s *FuncTestSuite) TestFuncCodeGen() {
	fw := primitives.FileWrapper{PackageName: "_func", FileName: "func_example.go"}

	funcGen := primitives.FuncGen{
		Name: "funcName",
	}

	fieldOne := primitives.Field{
		Name: "stringParam",
		Type: "string",
	}
	funcGen.AddField(fieldOne)

	returnField := primitives.Field{
		Name: "err",
		Type: "error",
	}
	funcGen.AddReturnField(returnField)

	resp := genFile(fw, funcGen)
	s.Assert().NotEmpty(resp)

	err := resp.Save(fw.FileName)
	s.Assert().Nil(err)
}
func TestFuncTestSuite(t *testing.T) {
	suite.Run(t, new(FuncTestSuite))
}
