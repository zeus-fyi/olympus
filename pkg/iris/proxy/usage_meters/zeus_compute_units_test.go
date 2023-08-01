package iris_usage_meters

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

var ctx = context.Background()

type IrisZeusComputeUnitsTestSuite struct {
	test_suites_base.TestSuite
}

func (s *IrisZeusComputeUnitsTestSuite) SetupTest() {

}

func (s *IrisZeusComputeUnitsTestSuite) TestRead() {
	p := make([]byte, 10)
	psm := &PayloadSizeMeter{R: bytes.NewReader(p)}
	n, err := psm.Read(p)
	s.NoError(err)
	s.Equal(10, n)
	s.Equal(int64(10), psm.N())
}

func (s *IrisZeusComputeUnitsTestSuite) TestAdd() {
	psm := &PayloadSizeMeter{}
	psm.Add(100)

	s.Equal(int64(100), psm.N())
}

func (s *IrisZeusComputeUnitsTestSuite) TestSizeInKB() {
	psm := &PayloadSizeMeter{}
	psm.Add(2048)

	s.Equal(2.0, psm.SizeInKB())
}

func (s *IrisZeusComputeUnitsTestSuite) TestZeusComputeUnitsConsumed() {
	// Payload less than 1KB
	psm := &PayloadSizeMeter{}
	psm.Add(500) // less than 1KB
	s.Equal(ZeusUnitsPerRequest+1, psm.ZeusComputeUnitsConsumed())

	// Payload exactly 1KB
	psm = &PayloadSizeMeter{}
	psm.Add(1024) // exactly 1KB
	s.Equal(ZeusUnitsPerRequest+1, psm.ZeusComputeUnitsConsumed())

	// Payload more than 1KB
	psm = &PayloadSizeMeter{}
	psm.Add(2048) // more than 1KB
	s.Equal(ZeusUnitsPerRequest+2, psm.ZeusComputeUnitsConsumed())
}

func TestIrisZeusComputeUnitsTestSuite(t *testing.T) {
	suite.Run(t, new(IrisZeusComputeUnitsTestSuite))
}
