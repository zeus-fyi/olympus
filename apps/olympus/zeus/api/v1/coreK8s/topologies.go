package coreK8s

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleTopologyActionRequest(c echo.Context) error {
	request := new(TopologyActionRequest)
	if err := c.Bind(request); err != nil {
		return err
	}

	return c.JSON(http.StatusBadRequest, nil)
}
