package artemis_validator_service_groups_models

import (
	"testing"

	"github.com/stretchr/testify/suite"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	artemis_client "github.com/zeus-fyi/zeus/pkg/artemis/client"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

type EthScheduledDeliveryTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *EthScheduledDeliveryTestSuite) SetupTest() {
	s.InitLocalConfigs()
	s.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
}

func (s *EthScheduledDeliveryTestSuite) TestInsert() {
	ArtemisClient = artemis_client.NewDefaultArtemisClient(s.Tc.ProductionLocalTemporalBearerToken)
	pubKey := "0x974C0c36265b7aa658b63A6121041AeE9e4DFd1b"
	//addr := common.HexToAddress(pubKey)
	//rr := artemis_req_types.SendEtherPayload{
	//	TransferArgs: artemis_req_types.TransferArgs{
	//		Amount:    big.NewInt(1).Mul(signing_automation_ethereum.Gwei, big.NewInt(int64(GweiThirtyTwoEth))),
	//		ToAddress: addr,
	//	},
	//}
	//rAddr := addr.String()
	//s.Require().Equal(pubKey, rAddr)
	//rx, err := ArtemisClient.SendEther(ctx, rr, artemis_client.ArtemisEthereumEphemeral)
	//s.Require().Nil(err)
	//s.Require().NotNil(rx)
	sd := artemis_autogen_bases.EthScheduledDelivery{
		DeliveryScheduleType: "networkReset",
		ProtocolNetworkID:    hestia_req_types.EthereumEphemeryProtocolNetworkID,
		Amount:               GweiThirtyTwoEth + GweiGasFees,
		Units:                "gwei",
		PublicKey:            pubKey,
	}
	err := InsertDeliverySchedule(ctx, sd)
	s.Require().Nil(err)
}
func TestEthScheduledDeliveryTestSuite(t *testing.T) {
	suite.Run(t, new(EthScheduledDeliveryTestSuite))
}
