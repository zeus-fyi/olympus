package create_infra

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/packages"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"k8s.io/apimachinery/pkg/util/rand"
)

type CreateBeaconInfraTestSuite struct {
	b hestia_test.BaseHestiaTestSuite
	conversions_test.ConversionsTestSuite
}

var LocalTemporalUserID = 7138958574876245567
var LocalTemporalOrgID = 7138983863666903883

func (s *CreateBeaconInfraTestSuite) TestInsertInfraBase() {
	p := structs.Path{
		PackageName: "",
		DirIn:       s.TestDirectory + "/mocks/demo",
		FnIn:        "deployment.yaml",
		DirOut:      s.TestDirectory + "/mocks/demo_out",
		FnOut:       "deployment.yaml",
		FilterFiles: string_utils.FilterOpts{DoesNotStartWithThese: []string{"cm-demo", "service"}},
	}
	err := s.Yr.ReadK8sWorkloadDir(p)
	s.Require().Nil(err)

	cw, err := s.Yr.CreateChartWorkloadFromNativeK8s()
	s.Require().Nil(err)

	pkg := packages.Packages{
		Chart: charts.Chart{
			ChartPackages: autogen_bases.ChartPackages{
				ChartPackageID:   0,
				ChartName:        rand.String(10),
				ChartVersion:     rand.String(10),
				ChartDescription: sql.NullString{},
			},
		},
		ChartWorkload: cw,
	}

	inf := NewCreateInfrastructure()
	inf.Packages = pkg
	ctx := context.Background()
	inf.Name = fmt.Sprintf("test_%d", s.Ts.UnixTimeStampNow())
	inf.OrgID, inf.UserID = LocalTemporalOrgID, LocalTemporalUserID

	err = inf.InsertInfraBase(ctx)
	s.Require().Nil(err)

	fmt.Println("ChartPackageID")
	fmt.Println(inf.ChartPackageID)
	fmt.Println("TopologyID")
	fmt.Println(inf.TopologyID)
}

func TestCreateBeaconInfraTestSuite(t *testing.T) {
	suite.Run(t, new(CreateBeaconInfraTestSuite))
}
