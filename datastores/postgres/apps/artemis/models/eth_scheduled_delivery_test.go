package artemis_validator_service_groups_models

import (
	"testing"

	"github.com/stretchr/testify/suite"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

type EthScheduledDeliveryTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *EthScheduledDeliveryTestSuite) SetupTest() {
	s.InitLocalConfigs()
	s.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
}

func (s *EthScheduledDeliveryTestSuite) TestInsert() {

	sd := artemis_autogen_bases.EthScheduledDelivery{
		DeliveryScheduleType: "networkReset",
		ProtocolNetworkID:    hestia_req_types.EthereumEphemeryProtocolNetworkID,
		Amount:               GweiThirtyTwoEth + GweiGasFees,
		Units:                "gwei",
		PublicKey:            "0x974C0c36265b7aa658b63A6121041AeE9e4DFd1b",
	}
	err := InsertDeliverySchedule(ctx, sd)
	s.Require().Nil(err)
}
func TestEthScheduledDeliveryTestSuite(t *testing.T) {
	suite.Run(t, new(EthScheduledDeliveryTestSuite))
}
