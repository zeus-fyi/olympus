package create_infra

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type TopologyCreateClassRequest struct {
	SkeletonBaseIDs []int `json:"skeletonBaseIDs,omitempty"`
	TopologyClassID int   `json:"topologyClassID,omitempty"`
}

type TopologyCreateClassResponse struct {
	ClassID int    `json:"classID"`
	Status  string `json:"status,omitempty"`
}

func (t *TopologyCreateClassRequest) CreateTopologyClass(c echo.Context) error {

	skeletonBaseID := c.FormValue("skeletonBaseID")
	t.SkeletonBaseIDs = string_utils.IntStringArrToIntArr(skeletonBaseID)

	topologyClassID := c.FormValue("topologyClassID")
	t.TopologyClassID = string_utils.IntStringParser(topologyClassID)

	//ctx := context.Background()

	//	log.Err(err).Interface("orgUser", ou).Msg("TopologyActionCreateRequest: CreateTopology, InsertInfraBase")
	//	return c.JSON(http.StatusInternalServerError, err)
	//}

	resp := TopologyCreateResponse{
		TopologyID:     0,
		SkeletonBaseID: 0,
	}
	return c.JSON(http.StatusOK, resp)
}
