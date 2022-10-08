package create

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/hera/code_templates/models"
)

type CreateStructNameExampleTestSuite struct {
	models.StructNameExampleTestSuite
}

func (s *CreateStructNameExampleTestSuite) TestQueryName() {
	ctx := context.Background()
	qp := models.CreateTestQueryNameParams()

	structExamples := StructNameExamples{}
	err := structExamples.StructNameExamplesFieldCase(ctx, qp)
	s.Require().Nil(err)
	s.Require().NotEmpty(structExamples)
}

func TestInsertStructNameExampleTestSuite(t *testing.T) {
	suite.Run(t, new(CreateStructNameExampleTestSuite))
}
