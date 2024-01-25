package zeus_v1_ai

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_stripe "github.com/zeus-fyi/olympus/pkg/hestia/stripe"
)

func AiSchemasHandler(c echo.Context) error {
	request := new(artemis_orchestrations.JsonSchemaDefinition)
	if err := c.Bind(request); err != nil {
		return err
	}
	if request == nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return CreateOrUpdateSchema(c, request)
}

func CreateOrUpdateSchema(c echo.Context, js *artemis_orchestrations.JsonSchemaDefinition) error {
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

	if js.SchemaStrID != "" {
		si, err := strconv.Atoi(js.SchemaStrID)
		if err != nil {
			log.Err(err).Msg("failed to parse int")
			return c.JSON(http.StatusBadRequest, nil)
		}
		js.SchemaID = si
	}
	for is, _ := range js.Fields {
		if js.Fields[is].FieldStrID != "" {
			sfi, err := strconv.Atoi(js.Fields[is].FieldStrID)
			if err != nil {
				log.Err(err).Msg("failed to parse int")
				return c.JSON(http.StatusBadRequest, nil)
			}
			js.Fields[is].FieldID = sfi
		}
	}
	err := artemis_orchestrations.CreateOrUpdateJsonSchema(c.Request().Context(), ou, js, nil)
	if err != nil {
		log.Err(err).Msg("failed to insert action")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, js)
}
