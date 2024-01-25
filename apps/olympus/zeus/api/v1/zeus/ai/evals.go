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
	err := ValidateEvalOps(&t.EvalFn)
	if err != nil {
		log.Err(err).Msg("failed to validate eval")
		return c.JSON(http.StatusBadRequest, nil)
	}
	err = artemis_orchestrations.InsertOrUpdateEvalFnWithMetrics(c.Request().Context(), ou, &t.EvalFn)
	if err != nil {
		log.Err(err).Msg("failed to insert evals")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, t.EvalFn)
}

func ValidateStrArrayPayload(em *artemis_orchestrations.EvalMetric) error {
	if em == nil || em.EvalMetricComparisonValues == nil {
		return nil
	}
	fv, err := strconv.ParseFloat(aws.StringValue(em.EvalMetricComparisonValues.EvalComparisonString), 64)
	if err != nil {
		log.Err(err).Msg("failed to parse float")
		return err
	}
	if fv < 0 {
		return errors.New("invalid value")
	}
	em.EvalMetricComparisonValues.EvalComparisonNumber = &fv
	em.EvalMetricComparisonValues.EvalComparisonString = nil
	return nil
}

func ValidateEvalOps(ef *artemis_orchestrations.EvalFn) error {
	if ef == nil {
		return nil
	}
	if ef.EvalStrID != nil {
		eid := *ef.EvalStrID
		eidInt, err := strconv.Atoi(eid)
		if err != nil {
			log.Err(err).Msg("failed to parse int")
			return err
		}
		ef.EvalID = aws.Int(eidInt)
	}

	for emi, em := range ef.Schemas {
		if len(em.SchemaStrID) > 0 {
			sid := em.SchemaStrID
			sidInt, err := strconv.Atoi(sid)
			if err != nil {
				log.Err(err).Msg("failed to parse int")
				return err
			}
			ef.Schemas[emi].SchemaID = sidInt
		}
		for fi, fe := range em.Fields {
			if len(fe.FieldStrID) > 0 {
				fid := fe.FieldStrID
				fidInt, err := strconv.Atoi(fid)
				if err != nil {
					log.Err(err).Msg("failed to parse int")
					return err
				}
				ef.Schemas[emi].Fields[fi].FieldID = fidInt
			}

			if len(fe.FieldName) <= 0 {
				return errors.New("invalid field name")
			}
			if len(fe.DataType) <= 0 {
				return errors.New("invalid field type")
			}
			for _, evm := range fe.EvalMetrics {
				err := ValidateEvalMetricOps(fe.DataType, evm)
				if err != nil {
					log.Err(err).Msg("failed to validate eval")
					return err
				}
			}
		}

	}
	return nil
}

func ValidateEvalMetricOps(dataType string, em *artemis_orchestrations.EvalMetric) error {
	if em == nil {
		return nil
	}
	switch dataType + "-" + em.EvalOperator {
	case "array[string]" + "-" + "length-less-than":
		err := ValidateStrArrayPayload(em)
		if err != nil {
			log.Err(err).Msg("failed to validate eval")
			return err
		}
	case "array[string]" + "-" + "length-less-than-eq":
		err := ValidateStrArrayPayload(em)
		if err != nil {
			log.Err(err).Msg("failed to validate eval")
			return err
		}
	case "array[string]" + "-" + "length-greater-than":
		err := ValidateStrArrayPayload(em)
		if err != nil {
			log.Err(err).Msg("failed to validate eval")
			return err
		}
	case "array[string]" + "-" + "length-greater-than-eq":
		err := ValidateStrArrayPayload(em)
		if err != nil {
			log.Err(err).Msg("failed to validate eval")
			return err
		}
	}
	return nil
}
