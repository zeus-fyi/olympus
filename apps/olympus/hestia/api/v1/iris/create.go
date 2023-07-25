package hestia_iris_v1_routes

import "github.com/labstack/echo/v4"

func CreateOrgRoutesRequestHandler(c echo.Context) error {
	request := new(CreateOrgRoutesRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.Create(c)
}

type CreateOrgRoutesRequest struct {
}

func (r *CreateOrgRoutesRequest) Create(c echo.Context) error {
	return nil
}
