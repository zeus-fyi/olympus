package iris_tdigest

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/caio/go-tdigest/v4"
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
	td, _ := tdigest.New()

	for i := 0; i < 10000; i++ {
		// Analogue to t.AddWeighted(rand.Float64(), 1)
		err := td.Add(rand.Float64())
		s.NoError(err)
	}
	fmt.Printf("p(.5) = %.6f\n", td.Quantile(0.5))
	fmt.Printf("CDF(Quantile(.5)) = %.6f\n", td.CDF(td.Quantile(0.5)))
}

func TestIrisTdigestTestSuite(t *testing.T) {
	suite.Run(t, new(IrisTdigestTestSuite))
}
