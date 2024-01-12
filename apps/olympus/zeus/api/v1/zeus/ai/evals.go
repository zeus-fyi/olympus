package zeus_v1_ai

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

type CreateOrUpdateEvalsRequest struct {
	artemis_orchestrations.EvalFn
}

func CreateOrUpdateEvalsRequestHandler(c echo.Context) error {
	request := new(CreateOrUpdateEvalsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateOrUpdateEval(c)
}

func (t *CreateOrUpdateEvalsRequest) CreateOrUpdateEval(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if ou.OrgID <= 0 || ou.UserID <= 0 {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	t.OrgID = ou.OrgID
	t.UserID = ou.UserID

	for _, em := range t.EvalFn.EvalMetrics {
		switch em.EvalMetricDataType + "-" + em.EvalOperator {
		case "array[string]" + "-" + "length-less-than":
			err := ValidateStrArrayPayload(&em)
			if err != nil {
				return c.JSON(http.StatusBadRequest, nil)
			}
		case "array[string]" + "-" + "length-less-than-eq":
			err := ValidateStrArrayPayload(&em)
			if err != nil {
				return c.JSON(http.StatusBadRequest, nil)
			}
		case "array[string]" + "-" + "length-greater-than":
			err := ValidateStrArrayPayload(&em)
			if err != nil {
				return c.JSON(http.StatusBadRequest, nil)
			}
		case "array[string]" + "-" + "length-greater-than-eq":
			err := ValidateStrArrayPayload(&em)
			if err != nil {
				return c.JSON(http.StatusBadRequest, nil)
			}
		}
	}
	err := artemis_orchestrations.InsertOrUpdateEvalFnWithMetrics(c.Request().Context(), ou, &t.EvalFn)
	if err != nil {
		log.Err(err).Msg("failed to insert evals")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}

func ValidateStrArrayPayload(em *artemis_orchestrations.EvalMetric) error {
	if em == nil {
		return nil
	}
	fv, err := strconv.ParseFloat(aws.StringValue(em.EvalComparisonString), 64)
	if err != nil {
		log.Err(err).Msg("failed to parse float")
		return err
	}

	if fv < 0 {
		return errors.New("invalid value")
	}
	em.EvalComparisonNumber = &fv
	em.EvalComparisonString = nil
	return nil
}
