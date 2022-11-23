package athena_routines

import "github.com/labstack/echo/v4"

func PauseRoutineHandler(c echo.Context) error {
	request := new(RoutineRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.PauseApp(c)
}

func ResumeRoutineHandler(c echo.Context) error {
	request := new(RoutineRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.ResumeApp(c)
}
