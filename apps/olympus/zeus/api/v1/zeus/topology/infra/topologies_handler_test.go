package infra

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
)

type TopologyActionRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
}

func TestTopologyActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyActionRequestTestSuite))
}
