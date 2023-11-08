package create_infra

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/zeus/zeus/system_config_drivers"
)

func CreateMatrixInfraActionRequestHandler(c echo.Context) error {
	request := new(CreateMatrixTopologyRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateMatrixDefinition(c)
}

type CreateMatrixTopologyRequest struct {
	system_config_drivers.MatrixDefinition
}

func (t *CreateMatrixTopologyRequest) CreateMatrixDefinition(c echo.Context) error {
	tmp := t.MatrixName
	fmt.Println(tmp)

	for _, _ = range t.Clusters {
		fmt.Println("cluster")
	}

	return c.JSON(http.StatusOK, "ok")
}
