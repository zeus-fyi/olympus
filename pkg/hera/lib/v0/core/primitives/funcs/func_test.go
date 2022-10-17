package funcs

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/fields"
)

type FuncTestSuite struct {
	suite.Suite
}

func (s *FuncTestSuite) TestFuncCodeGen() {
	fw := fields.FileWrapper{PackageName: "_func", FileName: "func_example.go"}

	funcGen := FuncGen{
		Name: "funcName",
	}

	fieldOne := fields.Field{
		Name: "stringParam",
		Type: "string",
	}
	funcGen.AddField(fieldOne)

	returnField := fields.Field{
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
