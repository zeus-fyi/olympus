package _struct

import (
	"testing"

	"github.com/stretchr/testify/suite"
	primitives2 "github.com/zeus-fyi/olympus/pkg/hera/codegen/cookbook/core/primitives"
	"github.com/zeus-fyi/olympus/pkg/hera/codegen/cookbook/core/template_test"
)

type StructTestSuite struct {
	template_test.TemplateTestSuite
}

func (s *StructTestSuite) TestCodeGen() {
	fw := primitives2.FileWrapper{PackageName: "_struct", FileName: "struct_example.go"}

	structToMake := primitives2.StructGen{
		Name:   "StructExample",
		Fields: nil,
	}
	fieldOne := primitives2.Field{
		Name: "IntField",
		Type: "int",
	}
	structToMake.AddField(fieldOne)

	fieldTwo := primitives2.Field{
		Name: "StringField",
		Type: "string",
	}
	structToMake.AddField(fieldTwo)

	err := genFile(fw, structToMake).Save(fw.FileName)
	s.Require().Nil(err)

	if s.Cleanup {
		s.DeleteFile(fw.FileName)
	}
}

func TestFuncTestSuite(t *testing.T) {
	suite.Run(t, new(StructTestSuite))
}
