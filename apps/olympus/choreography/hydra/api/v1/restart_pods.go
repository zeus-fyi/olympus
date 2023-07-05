package v1_hydra_choreography

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	pods_client "github.com/zeus-fyi/zeus/zeus/z_client/workloads/pods"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
	zeus_pods_reqs "github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types/pods"
)

func RestartPods(c echo.Context) error {
	par := zeus_pods_reqs.PodActionRequest{
		TopologyDeployRequest: zeus_req_types.TopologyDeployRequest{
			TopologyID: 0,
			CloudCtxNs: CloudCtxNs,
		},
	}
	pc := pods_client.NewPodsClientFromZeusClient(ZeusClient)
	resp, err := pc.DeletePods(context.Background(), par)
	if err != nil {
		log.Err(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.String(http.StatusOK, string(resp))
}
