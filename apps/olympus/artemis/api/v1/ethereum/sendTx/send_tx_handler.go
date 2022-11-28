package artemis_eth_txs

import "github.com/labstack/echo/v4"

//func SendEthTxHandler(c echo.Context) error {
//	request := new(EthereumTxSendRequest)
//	if err := c.Bind(request); err != nil {
//		return err
//	}
//	return request.SendTx(c)
//}

func SendSignedTxEthGoerliTxHandler(c echo.Context) error {
	request := new(EthereumSendSignedTxRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.SendGoerliSignedTx(c)
}

func SendEtherGoerliTxHandler(c echo.Context) error {
	request := new(EthereumSendEtherRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.SendEtherGoerliTx(c)
}
