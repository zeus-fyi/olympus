package read_infra

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_topologies "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies"
)

func ReadUsersTopologiesMetadata(c echo.Context) error {
	tr := read_topologies.NewReadTopologiesMetadataGroup()
	// from auth lookup
	ou := c.Get("orgUser").(org_users.OrgUser)
	ctx := context.Background()
	err := tr.SelectTopologiesMetadata(ctx, ou)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, tr.Slice)
}
