package poseidon_chain_snapshots

import (
	"github.com/labstack/echo/v4"
)

func UploadChainSnapshotHandler(c echo.Context) error {
	request := new(UploadChainSnapshotRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.Upload(c)
}
