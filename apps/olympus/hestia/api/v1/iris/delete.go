package hestia_iris_v1_routes

import "github.com/labstack/echo/v4"

func DeleteOrgRoutesRequestHandler(c echo.Context) error {
	request := new(DeleteOrgRoutesRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.Delete(c)
}

type DeleteOrgRoutesRequest struct {
}

func (r *DeleteOrgRoutesRequest) Delete(c echo.Context) error {
	return nil
}
