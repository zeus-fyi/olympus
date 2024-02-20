package zeus_v1_ai

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_stripe "github.com/zeus-fyi/olympus/pkg/hestia/stripe"
)

type CreateOrUpdateEntitiesRequest struct {
	artemis_entities.EntitiesFilter
}

func CreateOrUpdateEntitiesRequestHandler(c echo.Context) error {
	request := new(CreateOrUpdateEntitiesRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateOrUpdateEntity(c)
}

func (e *CreateOrUpdateEntitiesRequest) CreateOrUpdateEntity(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Info().Interface("ou", ou)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	isBillingSetup, err := hestia_stripe.DoesUserHaveBillingMethod(c.Request().Context(), ou.UserID)
	if err != nil {
		log.Error().Err(err).Msg("failed to check if user has billing method")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if !isBillingSetup {
		return c.JSON(http.StatusPreconditionFailed, nil)
	}
	mdb := artemis_entities.UserEntityMetadata{
		JsonData: e.MetadataJsonb,
		TextData: aws.String(e.MetadataText),
		Labels:   []artemis_entities.UserEntityMetadataLabel{},
	}
	var labels []artemis_entities.UserEntityMetadataLabel
	for _, lv := range e.Labels {
		labels = append(labels, artemis_entities.UserEntityMetadataLabel{
			Label: lv,
		})
	}
	mdb.Labels = labels
	urw := &artemis_entities.UserEntityWrapper{
		UserEntity: artemis_entities.UserEntity{
			Nickname:  e.Nickname,
			Platform:  e.Platform,
			FirstName: e.FirstName,
			LastName:  e.LastName,
			MdSlice: []artemis_entities.UserEntityMetadata{
				mdb,
			},
		},
		Ou: ou,
	}
	err = artemis_entities.InsertUserEntityLabeledMetadata(c.Request().Context(), urw)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusOK, urw)
}

type SelectEntitiesRequest struct {
	artemis_entities.EntitiesFilter
}

func SelectEntitiesRequestHandler(c echo.Context) error {
	request := new(SelectEntitiesRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.SelectEntities(c)
}

func (e *SelectEntitiesRequest) SelectEntities(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Info().Interface("ou", ou)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	isBillingSetup, err := hestia_stripe.DoesUserHaveBillingMethod(c.Request().Context(), ou.UserID)
	if err != nil {
		log.Error().Err(err).Msg("failed to check if user has billing method")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if !isBillingSetup {
		return c.JSON(http.StatusPreconditionFailed, nil)
	}
	evs, err := artemis_entities.SelectUserMetadataByNicknameAndPlatform(c.Request().Context(), e.Nickname, e.Platform, e.Labels, e.SinceUnixTimestamp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusOK, evs)
}