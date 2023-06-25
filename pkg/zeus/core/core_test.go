package zeus_core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CoreTestSuite struct {
	K8TestSuite
}

func (s *CoreTestSuite) TestK8Contexts() {
	kctx, err := s.K.GetContexts()
	s.Nil(err)
	s.Greater(len(kctx), 0)

	for k, ct := range kctx {
		fmt.Println(k, ct.Cluster)
	}
	s.K.SetContext("do-nyc1-do-nyc1-zeus-demo")
	fmt.Println(kctx)
}

func TestCoreTestSuite(t *testing.T) {
	suite.Run(t, new(CoreTestSuite))
}
