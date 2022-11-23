package athena_routines

import "github.com/labstack/echo/v4"

func HypnosKillReplaceRoutineHandler(c echo.Context) error {
	request := new(RoutineRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.HypnosKill(c)
}

func ResumeRoutineHandler(c echo.Context) error {
	request := new(RoutineRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.ResumeApp(c)
}

func SuspendRoutineHandler(c echo.Context) error {
	request := new(RoutineRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.SuspendApp(c)
}

func KillProcessRoutineHandler(c echo.Context) error {
	request := new(RoutineRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.KillProcess(c)
}
