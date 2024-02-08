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
	v1 "k8s.io/api/core/v1"
)

func DeployDeploymentHandlerWrapper(k autok8s_core.K8Util) func(c echo.Context) error {
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
		if request.Kns.TopologyBaseInfraWorkload.Deployment != nil {
			if request.Kns.CloudProvider == "ovh" && request.Kns.Context == "zeusfyi" && request.Kns.TopologyBaseInfraWorkload.Deployment.Spec.Template.Spec.ImagePullSecrets == nil {
				request.Kns.TopologyBaseInfraWorkload.Deployment.Spec.Template.Spec.ImagePullSecrets = []v1.LocalObjectReference{{
					Name: "zeus-fyi-ext",
				}}
			}
			if request.Kns.CloudCtxNs.Context != "do-sfo3-dev-do-sfo3-zeus" {
				if request.Kns.TopologyBaseInfraWorkload.Deployment.Spec.Template.Spec.Tolerations == nil {
					request.Kns.TopologyBaseInfraWorkload.Deployment.Spec.Template.Spec.Tolerations = []v1.Toleration{}
				}
				request.Kns.TopologyBaseInfraWorkload.Deployment.Spec.Template.Spec.Tolerations = []v1.Toleration{
					{
						Key:      fmt.Sprintf("org-%d", request.OrgUser.OrgID),
						Operator: "Equal",
						Value:    fmt.Sprintf("org-%d", request.OrgUser.OrgID),
						Effect:   "NoSchedule",
					},
				}
				if request.Kns.ClusterClassName != "" {
					request.Kns.TopologyBaseInfraWorkload.Deployment.Spec.Template.Spec.Tolerations = append(request.Kns.TopologyBaseInfraWorkload.Deployment.Spec.Template.Spec.Tolerations, v1.Toleration{
						Key:      "app",
						Operator: "Equal",
						Value:    request.Kns.ClusterClassName,
						Effect:   "NoSchedule",
					})
				}
			}
			log.Debug().Interface("kns", request.Kns).Msg("DeployDeploymentHandler: CreateDeploymentIfVersionLabelChangesOrDoesNotExist")
			_, err := k.CreateDeploymentIfVersionLabelChangesOrDoesNotExist(ctx, request.Kns.CloudCtxNs, request.Kns.TopologyBaseInfraWorkload.Deployment, nil)
			if err != nil {
				log.Err(err).Msg("DeployDeploymentHandler")
				return c.JSON(http.StatusInternalServerError, err)
			}
		} else {
			err := errors.New("no deployment workload was supplied")
			return c.JSON(http.StatusBadRequest, err)
		}
		return c.JSON(http.StatusOK, nil)
	}
}
