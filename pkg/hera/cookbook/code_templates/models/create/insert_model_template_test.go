package create

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/hera/cookbook/code_templates/models/test"
)

type CreateStructNameExampleTestSuite struct {
	test.StructNameExampleTestSuite
}

func (s *CreateStructNameExampleTestSuite) TestQueryName() {
	ctx := context.Background()
	qp := test.CreateTestQueryNameParams()

	structExamples := StructNameExamples{}
	err := structExamples.StructNameExamplesFieldCase(ctx, qp)
	s.Require().Nil(err)
	s.Require().NotEmpty(structExamples)
}

func TestInsertStructNameExampleTestSuite(t *testing.T) {
	suite.Run(t, new(CreateStructNameExampleTestSuite))
}
