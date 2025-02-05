package artemis_eth_txs

import (
	"time"

	"github.com/labstack/echo/v4"
)

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

func SendSignedTxEthEphemeralTxHandler(c echo.Context) error {
	request := new(EthereumSendSignedTxRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.SendEphemeralSignedTx(c)
}

func SendEtherEphemeralTxHandler(c echo.Context) error {
	request := new(EthereumSendEtherRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.SendEtherEphemeral(c)
}

func SendEtherEphemeralFaucetTxHandler(c echo.Context) error {
	//request := new(EthereumSendEtherRequest)
	//if err := c.Bind(request); err != nil {
	//	return err
	//}
	mw := NewLimiter(0, time.Duration(100)*time.Minute)
	return mw.ServeHTTP(c.Response().Writer, c.Request(), c)
}

func SendSignedTxEthMainnetTxHandler(c echo.Context) error {
	request := new(EthereumSendSignedTxRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return nil
}

func SendEtherMainnetTxHandler(c echo.Context) error {
	request := new(EthereumSendEtherRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return nil
}
