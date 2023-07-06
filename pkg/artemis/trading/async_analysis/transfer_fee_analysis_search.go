package async_analysis

import (
	"context"
	"fmt"

	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
)

func (c *ContractAnalysis) CalculateTransferFeeTaxRange(ctx context.Context) error {
	stepCountGwei := 2
	stepCount := 10
	for i := 1; i < stepCountGwei+1; i++ {
		size := i * 100
		unitValue := artemis_eth_units.GweiMultiple(size)
		feePerc, err := c.CalculateTransferFeeTax(ctx, unitValue)
		if err != nil {
			return err
		}

		fmt.Println(size, " gwei transfer: feePerc: ", feePerc.Numerator, feePerc.Denominator)
	}
	for i := 1; i < stepCount+1; i++ {
		unitValue := artemis_eth_units.EtherMultiple(i)
		feePerc, err := c.CalculateTransferFeeTax(ctx, unitValue)
		if err != nil {
			return err
		}
		fmt.Println(i, " eth transfer feePerc: ", feePerc.Numerator, feePerc.Denominator)
	}
	return nil
}
