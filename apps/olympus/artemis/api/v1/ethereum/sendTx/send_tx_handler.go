package artemis_eth_txs

import "github.com/labstack/echo/v4"

func SendEthTxHandler(c echo.Context) error {
	request := new(EthereumTxSendRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.SendTx(c)
}

func SendGoerliEthTxHandler(c echo.Context) error {
	request := new(EthereumTxSendRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.SendGoerliTx(c)
}
