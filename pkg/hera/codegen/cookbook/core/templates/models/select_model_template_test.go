package models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SelectStructNameExampleTestSuite struct {
	StructNameExampleTestSuite
}

func (s *SelectStructNameExampleTestSuite) TestSelectQueryName() {
	ctx := context.Background()
	qp := createTestQueryNameParams()

	structExamples := StructNameExamples{}
	err := structExamples.SelectStructNameExamplesFieldCase(ctx, qp)
	s.Require().Nil(err)
	s.Require().NotEmpty(structExamples)
}

func TestSelectStructNameExampleTestSuite(t *testing.T) {
	suite.Run(t, new(SelectStructNameExampleTestSuite))
}
