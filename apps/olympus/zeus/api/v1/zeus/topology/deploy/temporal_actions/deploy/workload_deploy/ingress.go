package internal_deploy

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/temporal_actions/base_request"
)

func DeployIngressHandlerWrapper(k autok8s_core.K8Util) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := context.Background()
		k8CfgInterface := c.Get("k8Cfg")
		if k8CfgInterface != nil {
			k8Cfg, ok := k8CfgInterface.(autok8s_core.K8Util) // Ensure the type assertion is correct
			if ok {
				k = k8Cfg
			}
		}
		// Attempt to retrieve the InternalDeploymentActionRequest from the context
		requestInterface := c.Get("internalDeploymentActionRequest")
		if requestInterface == nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal deployment action request not found in context"})
		}
		request, ok := requestInterface.(*base_request.InternalDeploymentActionRequest) // Ensure the type assertion is correct
		if !ok {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid request type"})
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
			if request.Kns.CloudCtxNs.CloudProvider == "ovh" && request.Kns.CloudCtxNs.Context == "zeusfyi" && request.Kns.Namespace == "flows" {
				ns := request.Kns.CloudCtxNs.Namespace
				if request.Kns.TopologyBaseInfraWorkload.Ingress.Spec.Rules != nil {
					for ind, _ := range request.Kns.TopologyBaseInfraWorkload.Ingress.Spec.Rules {
						request.Kns.TopologyBaseInfraWorkload.Ingress.Spec.Rules[ind].Host = fmt.Sprintf("api.%s.zeus.fyi", ns)
					}
				}
				if request.Kns.TopologyBaseInfraWorkload.Ingress.Spec.TLS != nil {
					for ind, _ := range request.Kns.TopologyBaseInfraWorkload.Ingress.Spec.TLS {
						for ind2, _ := range request.Kns.TopologyBaseInfraWorkload.Ingress.Spec.TLS[ind].Hosts {
							request.Kns.TopologyBaseInfraWorkload.Ingress.Spec.TLS[ind].Hosts[ind2] = fmt.Sprintf("api.%s.zeus.fyi", ns)
						}
						request.Kns.TopologyBaseInfraWorkload.Ingress.Spec.TLS[ind].SecretName = fmt.Sprintf("%s-api-tls", ns)
					}
				}
			}
			log.Debug().Interface("kns", request.Kns).Msg("DeployIngressHandler: CreateIngressIfVersionLabelChangesOrDoesNotExist")
			_, err := k.CreateIngressIfVersionLabelChangesOrDoesNotExist(ctx, request.Kns.CloudCtxNs, request.Kns.TopologyBaseInfraWorkload.Ingress, nil)
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
}
