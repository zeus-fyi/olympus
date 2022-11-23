package athena_routines

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/athena/pkg/routines"
)

type RoutineRequest struct {
	ClientName string `json:"clientName"`
}

func (t *RoutineRequest) ResumeApp(c echo.Context) error {
	ctx := context.Background()
	appName := routines.GetProcessName(t.ClientName)
	err := routines.KillProcessWithCtx(ctx, appName)
	if err != nil {
		log.Err(err).Msg("ResumeApp: KillProcessWithCtx")
		return c.JSON(http.StatusInternalServerError, err)
	}
	cmd := exec.Command("sh", "-c", "/scripts/start.sh")
	err = cmd.Run()
	if err != nil {
		log.Err(err).Msg("ResumeApp: Start Script")
		return err
	}
	return c.JSON(http.StatusOK, nil)
}

func (t *RoutineRequest) HypnosKill(c echo.Context) error {
	appName := routines.GetProcessName(t.ClientName)
	ctx := context.Background()
	err := routines.KillProcessWithCtx(ctx, appName)
	if err != nil {
		log.Err(err).Msg("HypnosKill: KillProcessWithCtx")
		return c.JSON(http.StatusInternalServerError, err)
	}
	err = t.InjectHypnos()
	if err != nil {
		log.Err(err).Msg("HypnosKill: InjectHypnos")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}

func (t *RoutineRequest) InjectHypnos() error {
	port := routines.GetProcessPorts(t.ClientName)
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("hypnos --port=%s", port))
	err := cmd.Run()
	if err != nil {
		log.Err(err).Msg("InjectHypnos")
		return err
	}
	return nil
}

func (t *RoutineRequest) KillProcess(c echo.Context) error {
	appName := routines.GetProcessName(t.ClientName)
	ctx := context.Background()
	err := routines.KillProcessWithCtx(ctx, appName)
	if err != nil {
		log.Err(err).Msg("KillProcess")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}

func (t *RoutineRequest) SuspendApp(c echo.Context) error {
	appName := routines.GetProcessName(t.ClientName)
	ctx := context.Background()
	err := routines.SuspendProcessWithCtx(ctx, appName)
	if err != nil {
		log.Err(err).Msg("SuspendApp")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}
