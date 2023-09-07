package hestia_iris_dashboard

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func TutorialToggleRequestHandler(c echo.Context) error {
	request := new(TutorialToggleRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.ToggleTutorialSetting(c)
}

type TutorialToggleRequest struct{}

func (t *TutorialToggleRequest) ToggleTutorialSetting(c echo.Context) error {
	// TODO
	fmt.Println("TODO")
	return c.JSON(http.StatusOK, nil)
}
