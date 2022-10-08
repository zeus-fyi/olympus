package read

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/hera/code_templates/models"
)

type ReadStructNameExampleTestSuite struct {
	models.StructNameExampleTestSuite
}

func (s *ReadStructNameExampleTestSuite) TestSelectQueryName() {
	ctx := context.Background()
	qp := models.CreateTestQueryNameParams()

	structExamples := StructNameExamples{}
	err := structExamples.StructNameExamplesFieldCase(ctx, qp)
	s.Require().Nil(err)
	s.Require().NotEmpty(structExamples)
}

func TestReadStructNameExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ReadStructNameExampleTestSuite))
}
