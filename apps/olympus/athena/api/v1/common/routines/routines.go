package athena_routines

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type RoutineRequest struct {
	ClientName string `json:"clientName"`
}

func (t *RoutineRequest) ResumeApp(c echo.Context) error {
	err := t.terminateApp("hypnos")
	if err != nil {
		log.Err(err).Msg("TerminateApp")
		return c.JSON(http.StatusInternalServerError, err)
	}
	_, err = exec.Command("sh", "-c", "/scripts/start.sh").Output()
	if err != nil {
		log.Err(err).Msg("ResumeApp")
		return err
	}
	return c.JSON(http.StatusOK, nil)
}

func (t *RoutineRequest) PauseApp(c echo.Context) error {
	appName := ""
	switch t.ClientName {
	case "lighthouse":
		appName = "lighthouse"
	case "geth":
		appName = "geth"
	}
	err := t.terminateApp(appName)
	if err != nil {
		log.Err(err).Msg("TerminateApp")
		return c.JSON(http.StatusInternalServerError, err)
	}
	err = t.InjectHypnos()
	if err != nil {
		log.Err(err).Msg("InjectHypnos")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}

func (t *RoutineRequest) terminateApp(appName string) error {
	_, err := exec.Command("sh", "-c", fmt.Sprintf("pkill -SIGINT %s", appName)).Output()
	if err != nil {
		return err
	}
	return nil
}

func (t *RoutineRequest) InjectHypnos() error {
	port := ""
	switch t.ClientName {
	case "lighthouse":
		port = "5052"
	case "geth":
		port = "8545"
	}
	_, err := exec.Command("sh", "-c", fmt.Sprintf("hypnos --port=%s", port)).Output()
	if err != nil {
		return err
	}
	return nil
}
