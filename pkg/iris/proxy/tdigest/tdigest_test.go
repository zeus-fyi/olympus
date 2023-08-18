package iris_tdigest

import (
	"context"
	"log"
	"testing"

	"github.com/influxdata/tdigest"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

var ctx = context.Background()

type IrisTdigestTestSuite struct {
	test_suites_base.TestSuite
}

func (s *IrisTdigestTestSuite) SetupTest() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
}

func (s *IrisTdigestTestSuite) TestTdigest() {
	td := tdigest.NewWithCompression(100)

	for i := 0; i < 10000; i++ {
		td.Add(float64(i), 1)
	}
	// Compute Quantiles
	log.Println("50th", td.Quantile(0.5))
	log.Println("75th", td.Quantile(0.75))
	log.Println("90th", td.Quantile(0.9))
	log.Println("99th", td.Quantile(0.99))

	// Compute CDFs
	log.Println("CDF(1) = ", td.CDF(1))
	log.Println("CDF(2) = ", td.CDF(2))
	log.Println("CDF(3) = ", td.CDF(3))
	log.Println("CDF(4) = ", td.CDF(4))
	log.Println("CDF(5) = ", td.CDF(5))

}

func TestIrisTdigestTestSuite(t *testing.T) {
	suite.Run(t, new(IrisTdigestTestSuite))
}
