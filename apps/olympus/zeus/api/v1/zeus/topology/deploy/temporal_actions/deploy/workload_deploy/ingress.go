package internal_deploy

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/temporal_actions/base_request"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

func DeployIngressHandler(c echo.Context) error {
	ctx := context.Background()
	request := new(base_request.InternalDeploymentActionRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if request.Kns.TopologyBaseInfraWorkload.Ingress != nil {
		if request.Kns.CloudCtxNs.Context != "do-sfo3-dev-do-sfo3-zeus" {
			ns := request.Kns.CloudCtxNs.Namespace
			if request.Kns.TopologyBaseInfraWorkload.Ingress.Spec.Rules != nil {
				for ind, _ := range request.Kns.TopologyBaseInfraWorkload.Ingress.Spec.Rules {
					request.Kns.TopologyBaseInfraWorkload.Ingress.Spec.Rules[ind].Host = fmt.Sprintf("%s.zeus.fyi", ns)
				}
			}
			if request.Kns.TopologyBaseInfraWorkload.Ingress.Spec.TLS != nil {
				for ind, _ := range request.Kns.TopologyBaseInfraWorkload.Ingress.Spec.TLS {
					for ind2, _ := range request.Kns.TopologyBaseInfraWorkload.Ingress.Spec.TLS[ind].Hosts {
						request.Kns.TopologyBaseInfraWorkload.Ingress.Spec.TLS[ind].Hosts[ind2] = fmt.Sprintf("%s.zeus.fyi", ns)
					}
					request.Kns.TopologyBaseInfraWorkload.Ingress.Spec.TLS[ind].SecretName = fmt.Sprintf("%s-tls", ns)
				}
			}
		}
		log.Debug().Interface("kns", request.Kns).Msg("DeployIngressHandler: CreateIngressIfVersionLabelChangesOrDoesNotExist")
		_, err := zeus.K8Util.CreateIngressIfVersionLabelChangesOrDoesNotExist(ctx, request.Kns.CloudCtxNs, request.Kns.TopologyBaseInfraWorkload.Ingress, nil)
		if err != nil {
			log.Err(err).Msg("DeployIngressHandler")
			return c.JSON(http.StatusInternalServerError, err)
		}
	} else {
		err := errors.New("no ingress workload was supplied")
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, nil)
}
