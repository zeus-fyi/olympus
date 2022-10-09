package _struct

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/template_test"
)

type StructTestSuite struct {
	template_test.TemplateTestSuite
}

func (s *StructTestSuite) TestCodeGen() {
	fw := primitives.FileWrapper{PackageName: "_struct", FileName: "struct_example.go"}

	structToMake := primitives.StructGen{
		Name:   "StructExample",
		Fields: nil,
	}
	fieldOne := primitives.Field{
		Name: "IntField",
		Type: "int",
	}
	structToMake.AddField(fieldOne)

	fieldTwo := primitives.Field{
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
