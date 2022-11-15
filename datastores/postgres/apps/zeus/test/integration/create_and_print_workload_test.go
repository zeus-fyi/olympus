package integration

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/packages"
	create_infra "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/definitions/classes/bases/infra"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"k8s.io/apimachinery/pkg/util/rand"
)

type IntegrationTestSuite struct {
	b hestia_test.BaseHestiaTestSuite
	conversions_test.ConversionsTestSuite
}

var LocalTemporalUserID = 7138958574876245567
var LocalTemporalOrgID = 7138983863666903883

func (s *IntegrationTestSuite) TestInsertInfraBase() {
	s.ChangeToTestDirectory()

	p := structs.Path{
		PackageName: "",
		DirIn:       s.TestDirectory + "/mocks/consensus_client",
		FnIn:        "statefulset.yaml",
		DirOut:      s.TestDirectory + "/mocks/consensus_client_out",
		FnOut:       "statefulset.yaml",
		FilterFiles: string_utils.FilterOpts{DoesNotStartWithThese: []string{"cm-lighthouse", "service"}},
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

	inf := create_infra.NewCreateInfrastructure()
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

	tr := read_topology.NewInfraTopologyReader()

	tr.TopologyID = inf.TopologyID
	tr.OrgID = LocalTemporalOrgID
	tr.UserID = LocalTemporalUserID
	err = tr.SelectTopology(ctx)
	s.Require().Nil(err)

	chart := tr.Chart

	b, err := json.Marshal(chart.StatefulSet.K8sStatefulSet)
	s.Require().Nil(err)

	err = s.Yr.WriteYamlConfig(p, b)
	s.Require().Nil(err)

}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
