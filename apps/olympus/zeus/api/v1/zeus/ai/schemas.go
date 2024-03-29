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

	if js == nil {
		return c.JSON(http.StatusBadRequest, nil)
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

type JsonSchemaDefinitionsReader struct {
}

func AiSchemaHandler(c echo.Context) error {
	request := new(JsonSchemaDefinitionsReader)
	if err := c.Bind(request); err != nil {
		return err
	}
	if request == nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return GetSchemas(c)
}

func GetSchemas(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	//isBillingSetup, berr := hestia_stripe.DoesUserHaveBillingMethod(c.Request().Context(), ou.UserID)
	//if berr != nil {
	//	log.Error().Err(berr).Msg("failed to check if user has billing method")
	//	return c.JSON(http.StatusInternalServerError, nil)
	//}
	//if !isBillingSetup {
	//	return c.JSON(http.StatusPreconditionFailed, nil)
	//}
	jsds, err := artemis_orchestrations.SelectJsonSchemaByOrg(c.Request().Context(), ou)
	if err != nil {
		log.Err(err).Msg("failed to select schemas")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, jsds)
}

//type GetSchemaRequest struct {
//}
//
//func GetSchemaRequestHandler(c echo.Context) error {
//	request := new(GetSchemaRequest)
//	if err := c.Bind(request); err != nil {
//		return err
//	}
//	// Extracting the :id parameter from the route
//	idParam := c.Param("id")
//	id, err := strconv.Atoi(idParam)
//	if err != nil {
//		log.Err(err).Msg("invalid ID parameter")
//		return c.JSON(http.StatusBadRequest, "invalid ID parameter")
//	}
//	return request.GetSchemaByID(c, id)
//}
//
//func (t *GetSchemaRequest) GetSchemaByID(c echo.Context, id int) error {
//	ou, ok := c.Get("orgUser").(org_users.OrgUser)
//	if !ok {
//		return c.JSON(http.StatusInternalServerError, nil)
//	}
//	if ou.OrgID <= 0 || ou.UserID <= 0 {
//		return c.JSON(http.StatusInternalServerError, nil)
//	}
//
//	ret, err := artemis_orchestrations.SelectJsonSchemaByOrg(c.Request().Context(), ou, id)
//	if err != nil {
//		log.Err(err).Msg("failed to get retrievals")
//		return c.JSON(http.StatusInternalServerError, nil)
//	}
//	return c.JSON(http.StatusOK, ret)
//}
