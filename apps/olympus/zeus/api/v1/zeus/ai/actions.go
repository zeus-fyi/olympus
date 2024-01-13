package zeus_v1_ai

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_stripe "github.com/zeus-fyi/olympus/pkg/hestia/stripe"
)

func AiActionsHandler(c echo.Context) error {
	request := new(artemis_orchestrations.TriggerAction)
	if err := c.Bind(request); err != nil {
		return err
	}
	if request == nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return CreateOrUpdateAction(c, request)
}

func CreateOrUpdateAction(c echo.Context, act *artemis_orchestrations.TriggerAction) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	isBillingSetup, berr := hestia_stripe.DoesUserHaveBillingMethod(c.Request().Context(), ou.UserID)
	if berr != nil {
		log.Error().Err(berr).Msg("failed to check if user has billing method")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if !isBillingSetup {
		return c.JSON(http.StatusPreconditionFailed, nil)
	}

	act.EvalTriggerActions = append(act.EvalTriggerActions, act.EvalTriggerAction)
	err := artemis_orchestrations.CreateOrUpdateTriggerAction(c.Request().Context(), ou, act)
	if err != nil {
		log.Err(err).Msg("failed to insert action")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, act)
}

type ActionApprovalRequest struct {
	RequestedState         string                                        `json:"requestedState"`
	TriggerActionsApproval artemis_orchestrations.TriggerActionsApproval `json:"triggerApproval"`
}

func AiActionsApprovalHandler(c echo.Context) error {
	request := new(ActionApprovalRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	if request == nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return UpdateActionApproval(c, request)
}

func UpdateActionApproval(c echo.Context, act *ActionApprovalRequest) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	isBillingSetup, berr := hestia_stripe.DoesUserHaveBillingMethod(c.Request().Context(), ou.UserID)
	if berr != nil {
		log.Error().Err(berr).Msg("failed to check if user has billing method")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if !isBillingSetup {
		return c.JSON(http.StatusPreconditionFailed, nil)
	}

	aptr := &act.TriggerActionsApproval
	aptr.ApprovalState = act.RequestedState
	err := artemis_orchestrations.CreateOrUpdateTriggerActionApproval(c.Request().Context(), ou, aptr)
	if err != nil {
		log.Err(err).Msg("failed to insert action")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, aptr)
}
