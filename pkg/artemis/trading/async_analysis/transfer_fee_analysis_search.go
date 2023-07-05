package async_analysis

import (
	"context"
	"fmt"

	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/units"
)

func (c *ContractAnalysis) CalculateTransferFeeTaxRange(ctx context.Context) error {
	stepCount := 10
	for i := 0; i < stepCount; i++ {
		unitValue := artemis_eth_units.EtherMultiple(i)
		tft, err := c.CalculateTransferFeeTax(ctx, unitValue)
		if err != nil {
			return err
		}
		fmt.Printf("TransferFeeTax: %s\n", tft.Quotient().String())
	}
	return nil
}
