package _func

import (
	"testing"

	"github.com/stretchr/testify/suite"
	primitives2 "github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives"
)

type FuncTestSuite struct {
	suite.Suite
}

func (s *FuncTestSuite) TestFuncCodeGen() {
	fw := primitives2.FileWrapper{PackageName: "_func", FileName: "func_example.go"}

	funcGen := primitives2.FuncGen{
		Name: "funcName",
	}

	fieldOne := primitives2.Field{
		Name: "stringParam",
		Type: "string",
	}
	funcGen.AddField(fieldOne)

	returnField := primitives2.Field{
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
