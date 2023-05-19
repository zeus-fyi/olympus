package artemis_validator_service_groups_models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

type MevTxTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *MevTxTestSuite) TestInsertNodes() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	p := filepaths.Path{
		PackageName: "",
		DirIn:       "/Users/alex/go/Olympus/olympus/test/p2p",
		DirOut:      "",
		FnIn:        "all-nodes.json",
		FnOut:       "",
		Env:         "",
		FilterFiles: nil,
	}
	b := p.ReadFileInPath()
	var nodes P2PNodes
	err := json.Unmarshal(b, &nodes)
	s.Require().Nil(err)
	err = InsertP2PNodes(ctx, artemis_autogen_bases.EthP2PNodes{
		ID:                0,
		ProtocolNetworkID: 0,
		Nodes:             string(b),
	})
	s.Require().Nil(err)
	selectedNodes, err := SelectP2PNodes(ctx, 0)
	s.Require().Nil(err)
	s.Require().Equal(len(nodes), len(selectedNodes))
}

func TestMevTxTestSuite(t *testing.T) {
	suite.Run(t, new(MevTxTestSuite))
}
