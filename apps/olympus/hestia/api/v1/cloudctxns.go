package v1hestia

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	create_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/topology"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type CreateTopologiesOrgCloudCtxNsRequest struct {
	OrgID int `db:"org_id" json:"orgID"`
	zeus_common_types.CloudCtxNs
}

func CreateTopologiesOrgCloudCtxNsRequestHandler(c echo.Context) error {
	request := new(CreateTopologiesOrgCloudCtxNsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateTopologiesOrgCloudCtxNsRequest(c)
}

func (t *CreateTopologiesOrgCloudCtxNsRequest) CreateTopologiesOrgCloudCtxNsRequest(c echo.Context) error {
	ctx := context.Background()

	orgCloudReq := create_topology.NewCreateTopologiesOrgCloudCtxNs(t.OrgID, t.CloudCtxNs)
	err := orgCloudReq.InsertTopologyAccessCloudCtxNs(ctx, t.OrgID, t.CloudCtxNs)
	if err != nil {
		log.Err(err).Interface("orgCloudReq", orgCloudReq).Msg("CreateTopologiesOrgCloudCtxNsRequest error")
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusOK, orgCloudReq)
}
