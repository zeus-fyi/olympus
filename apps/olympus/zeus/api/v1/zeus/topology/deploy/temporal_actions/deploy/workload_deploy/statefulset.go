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

func DeployStatefulSetHandler(c echo.Context) error {
	ctx := context.Background()
	request := new(base_request.InternalDeploymentActionRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if request.StatefulSet != nil {
		if request.Kns.CloudCtxNs.Context != "do-sfo3-dev-do-sfo3-zeus" {
			request.StatefulSet.Spec.Template.Spec.Tolerations = []v1.Toleration{
				{
					Key:      fmt.Sprintf("org-%d", request.OrgUser.OrgID),
					Operator: "Equal",
					Value:    fmt.Sprintf("org-%d", request.OrgUser.OrgID),
					Effect:   "NoSchedule",
				},
			}
		}
		log.Debug().Interface("kns", request.Kns).Msg("DeployStatefulSetHandler: CreateStatefulSetIfVersionLabelChangesOrDoesNotExist")
		_, err := zeus.K8Util.CreateStatefulSetIfVersionLabelChangesOrDoesNotExist(ctx, request.Kns.CloudCtxNs, request.StatefulSet, nil)
		if err != nil {
			log.Err(err).Msg("DeployStatefulSetHandler")
			return c.JSON(http.StatusInternalServerError, err)
		}
	} else {
		err := errors.New("no statefulset workload was supplied")
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, nil)
}
