package async_analysis

import (
	"context"

	"github.com/rs/zerolog/log"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
)

// FindERC20BalanceOfSlotNumber finds the erc20 token balance of
func (c *ContractAnalysis) FindERC20BalanceOfSlotNumber(ctx context.Context) error {
	err := c.u.Web3Client.SetERC20BalanceBruteForce(ctx, c.SmartContractAddr, c.UserA.Address().String(), artemis_eth_units.EtherMultiple(100))
	if err != nil {
		log.Err(err).Msg("ContractAnalysis: SetERC20BalanceBruteForce")
		return err
	}
	return nil
}

//// FindERC20BalanceOfSlotNumber2 finds the erc20 token balance of
//func (c *ContractAnalysis) FindERC20BalanceOfSlotNumber2(ctx context.Context, af *abi.ABI) error {
//	err := c.u.Web3Client.SetERC20BalanceBruteForceCustomErc20(ctx, c.SmartContractAddr, c.UserA.Address().String(), af, artemis_eth_units.EtherMultiple(100))
//	if err != nil {
//		log.Err(err).Msg("ContractAnalysis: SetERC20BalanceBruteForceCustomErc20")
//		return err
//	}
//	return nil
//}
