package artemis_eth_txs

import "github.com/labstack/echo/v4"

type EthereumTxSendRequest struct {
	// TODO payload
}

func (t *EthereumTxSendRequest) SendTx(c echo.Context) error {
	// TODO call artemis
	return nil
}

func (t *EthereumTxSendRequest) SendGoerliTx(c echo.Context) error {

	// TODO call artemis
	return nil
}
