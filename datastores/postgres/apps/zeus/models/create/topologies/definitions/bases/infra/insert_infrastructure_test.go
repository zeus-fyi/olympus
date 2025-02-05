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
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"k8s.io/apimachinery/pkg/util/rand"
)

type CreateInfraTestSuite struct {
	b hestia_test.BaseHestiaTestSuite
	conversions_test.ConversionsTestSuite
}

func (s *CreateInfraTestSuite) TestInsertInfraBase() {
	p := filepaths.Path{
		PackageName: "",
		DirIn:       s.TestDirectory + "/mocks/demo",
		FnIn:        "statefulset.yaml",
		DirOut:      s.TestDirectory + "/tempout",
		FnOut:       "statefulset.yaml",
		FilterFiles: string_utils.FilterOpts{DoesNotStartWithThese: []string{"service"}},
	}
	err := s.Yr.ReadK8sWorkloadDir(p)
	s.Require().Nil(err)

	cw, err := s.Yr.CreateChartWorkloadFromTopologyBaseInfraWorkload()
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
	//inf.OrgID, inf.UserID = s.b.NewTestOrgAndUser()
	inf.OrgID = s.Tc.ProductionLocalTemporalOrgID
	inf.UserID = s.Tc.ProductionLocalTemporalUserID
	fmt.Println("OrgID")
	fmt.Println(inf.OrgID)

	fmt.Println("UserID")
	fmt.Println(inf.UserID)

	err = inf.InsertInfraBase(ctx)
	s.Require().Nil(err)

	fmt.Println("ChartPackageID")
	fmt.Println(inf.ChartPackageID)
	fmt.Println("TopologyID")
	fmt.Println(inf.TopologyID)
}

func TestCreateInfraTestSuite(t *testing.T) {
	suite.Run(t, new(CreateInfraTestSuite))
}
