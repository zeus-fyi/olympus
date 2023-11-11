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
	v1 "k8s.io/api/core/v1"
)

func DeployDeploymentHandler(c echo.Context) error {
	ctx := context.Background()
	request := new(base_request.InternalDeploymentActionRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
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
		_, err := zeus.K8Util.CreateDeploymentIfVersionLabelChangesOrDoesNotExist(ctx, request.Kns.CloudCtxNs, request.Kns.TopologyBaseInfraWorkload.Deployment, nil)
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
