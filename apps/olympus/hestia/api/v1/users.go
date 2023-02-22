package v1hestia

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	create_org_users "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/org_users"
)

type CreateUserRequest struct {
	Metadata string `db:"metadata" json:"metadata"`
}

func CreateUserHandler(c echo.Context) error {
	request := new(CreateUserRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateUser(c)
}

func (uc *CreateUserRequest) CreateUser(c echo.Context) error {
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	u := create_org_users.NewCreateOrgUserWithOrgID(ou.OrgID)
	b, err := json.Marshal(uc.Metadata)
	if err != nil {
		log.Err(err).Interface("user", u).Msg("CreateUserRequest, CreateUser error")
		return c.JSON(http.StatusBadRequest, b)
	}
	err = u.InsertOrgUser(ctx, b)
	if err != nil {
		log.Err(err).Interface("user", u).Msg("CreateUserRequest, CreateUser error")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, u)
}

const (
	DemoUsersCreateRoute = "/users/demo/create"
)

type CreateDemoUserRequest struct {
	KeyName   string `json:"keyName"`
	Metadata  string `json:"metadata"`
	ServiceID int    `json:"serviceID"`
}

func CreateDemoUserHandler(c echo.Context) error {
	request := new(CreateDemoUserRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateUserWithKeyServiceDemo(c)
}

func (uc *CreateDemoUserRequest) CreateUserWithKeyServiceDemo(c echo.Context) error {
	ctx := context.Background()
	ou := create_org_users.OrgUser{}
	b, err := json.Marshal(uc.Metadata)
	if err != nil {
		log.Err(err).Interface("userMetadata", uc.Metadata).Msg("CreateUserRequest: marshal metadata error")
		return c.JSON(http.StatusBadRequest, b)
	}
	key, err := ou.InsertDemoOrgUserWithNewKey(ctx, b, uc.KeyName, uc.ServiceID)
	if err != nil {
		log.Err(err).Interface("userMetadata", uc.Metadata).Msg("CreateUserWithKeyServiceDemo, insert error")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	resp := make(map[string]string)
	resp["bearer"] = key
	return c.JSON(http.StatusOK, resp)
}
