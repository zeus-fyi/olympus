package zeus_v1_ai

import (
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

type CreateOrUpdateRetrievalRequest struct {
	artemis_orchestrations.RetrievalItem
}

func CreateOrUpdateRetrievalRequestHandler(c echo.Context) error {
	request := new(CreateOrUpdateRetrievalRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateOrUpdateRetrieval(c)
}

func (t *CreateOrUpdateRetrievalRequest) CreateOrUpdateRetrieval(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if ou.OrgID <= 0 || ou.UserID <= 0 {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if t.RetrievalName == "" || t.RetrievalPlatform == "" || (aws.StringValue(t.RetrievalKeywords) == "" && aws.StringValue(t.RetrievalPrompt) == "" && t.RetrievalGroup == "") {
		return c.JSON(http.StatusBadRequest, nil)
	}
	if t.RetrievalStrID != nil {
		rid, err := strconv.Atoi(*t.RetrievalStrID)
		if err != nil {
			log.Err(err).Msg("failed to parse int")
			return c.JSON(http.StatusBadRequest, nil)
		}
		t.RetrievalID = aws.Int(rid)
	}
	err := artemis_orchestrations.InsertRetrieval(c.Request().Context(), ou, &t.RetrievalItem)
	if err != nil {
		log.Err(err).Msg("failed to insert ret")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, t.RetrievalItem)
}
