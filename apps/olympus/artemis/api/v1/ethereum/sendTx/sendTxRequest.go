package sendEthTx

import "github.com/labstack/echo/v4"

type EthereumTxSendRequest struct {
}

func (t *EthereumTxSendRequest) SendTx(c echo.Context) error {

	// TODO call artemis
	return nil
}
