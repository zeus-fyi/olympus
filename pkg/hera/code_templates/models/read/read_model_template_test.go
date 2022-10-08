package read

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/hera/code_templates/models/test"
)

type ReadStructNameExampleTestSuite struct {
	test.StructNameExampleTestSuite
}

func (s *ReadStructNameExampleTestSuite) TestSelectQueryName() {
	ctx := context.Background()
	qp := test.CreateTestQueryNameParams()

	structExamples := StructNameExamples{}
	err := structExamples.StructNameExamplesFieldCase(ctx, qp)
	s.Require().Nil(err)
	s.Require().NotEmpty(structExamples)
}

func TestReadStructNameExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ReadStructNameExampleTestSuite))
}
