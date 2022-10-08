package models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type InsertStructNameExampleTestSuite struct {
	StructNameExampleTestSuite
}

func (s *InsertStructNameExampleTestSuite) TestQueryName() {
	ctx := context.Background()
	qp := createTestQueryNameParams()

	structExamples := StructNameExamples{}
	err := structExamples.InsertStructNameExamplesFieldCase(ctx, qp)
	s.Require().Nil(err)
	s.Require().NotEmpty(structExamples)
}

func TestInsertStructNameExampleTestSuite(t *testing.T) {
	suite.Run(t, new(InsertStructNameExampleTestSuite))
}
