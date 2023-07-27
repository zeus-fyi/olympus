package hestia_iris_v1_routes

import "github.com/labstack/echo/v4"

func UpdateOrgRoutesRequestHandler(c echo.Context) error {
	request := new(UpdateOrgRoutesRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.Update(c)
}

type UpdateOrgRoutesRequest struct {
}

func (r *UpdateOrgRoutesRequest) Update(c echo.Context) error {
	return nil
}
