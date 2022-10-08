package models

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type StructNameExampleTestSuite struct {
	test_suites.PGTestSuite
}

func TestStructNameExampleTestSuite(t *testing.T) {
	suite.Run(t, new(StructNameExampleTestSuite))
}
