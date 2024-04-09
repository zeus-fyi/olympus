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

func CreateOrUpdateEntitiesRequestHandler(c echo.Context) error {
	var entitiesFilters []artemis_entities.EntitiesFilter
	if err := c.Bind(&entitiesFilters); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return CreateOrUpdateEntities(c, entitiesFilters)
}

func CreateOrUpdateEntities(c echo.Context, ef []artemis_entities.EntitiesFilter) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Info().Interface("ou", ou)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	//isBillingSetup, err := hestia_stripe.DoesUserHaveBillingMethod(c.Request().Context(), ou.UserID)
	//if err != nil {
	//	log.Error().Err(err).Msg("failed to check if user has billing method")
	//	return c.JSON(http.StatusInternalServerError, nil)
	//}
	//if !isBillingSetup {
	//	return c.JSON(http.StatusPreconditionFailed, nil)
	//}
	for _, e := range ef {
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
		err := artemis_entities.InsertUserEntityLabeledMetadata(c.Request().Context(), urw)
		if err != nil {
			log.Err(err).Msg("failed to insert user entity")
			return c.JSON(http.StatusBadRequest, nil)
		}
	}
	return c.JSON(http.StatusOK, nil)
}

type CreateOrUpdateEntityRequest struct {
	artemis_entities.EntitiesFilter
}

func CreateOrUpdateEntityRequestHandler(c echo.Context) error {
	request := new(CreateOrUpdateEntityRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateOrUpdateEntity(c)
}

func (e *CreateOrUpdateEntityRequest) CreateOrUpdateEntity(c echo.Context) error {
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

type SearchEntitiesRequest struct {
	artemis_entities.EntitiesFilter
}

func SelectEntitiesRequestHandler(c echo.Context) error {
	request := new(SearchEntitiesRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.SelectEntities(c)
}

func (e *SearchEntitiesRequest) SelectEntities(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Info().Interface("ou", ou)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	//isBillingSetup, err := hestia_stripe.DoesUserHaveBillingMethod(c.Request().Context(), ou.UserID)
	//if err != nil {
	//	log.Error().Err(err).Msg("failed to check if user has billing method")
	//	return c.JSON(http.StatusInternalServerError, nil)
	//}
	//if !isBillingSetup {
	//	return c.JSON(http.StatusPreconditionFailed, nil)
	//}
	evs, err := artemis_entities.SelectUserMetadataByProvidedFields(c.Request().Context(), ou, e.Nickname, e.Platform, e.Labels, e.SinceUnixTimestamp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusOK, evs)
}
