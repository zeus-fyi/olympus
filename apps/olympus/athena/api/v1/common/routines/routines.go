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

func (t *RoutineRequest) Start(c echo.Context) error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "func", "Start")

	cmd := exec.Command("sh", "-c", "/scripts/start.sh")
	err := cmd.Run()
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("Start: Start Script")
		return err
	}
	return c.JSON(http.StatusOK, nil)
}

func (t *RoutineRequest) HypnosKill(c echo.Context) error {
	appName := routines.GetProcessName(t.ClientName)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "func", "HypnosKill")
	err := routines.KillProcessWithCtx(ctx, appName)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("HypnosKill: KillProcessWithCtx")
		return c.JSON(http.StatusInternalServerError, err)
	}
	err = t.InjectHypnos()
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("HypnosKill: InjectHypnos")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}

func (t *RoutineRequest) InjectHypnos() error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "func", "InjectHypnos")
	port := routines.GetProcessPorts(t.ClientName)
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("hypnos --port=%s", port))
	err := cmd.Run()
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("InjectHypnos")
		return err
	}
	return nil
}

func (t *RoutineRequest) Kill(c echo.Context) error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "func", "Kill")
	appName := routines.GetProcessName(t.ClientName)
	err := routines.KillProcessWithCtx(ctx, appName)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("Kill")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}

func (t *RoutineRequest) Suspend(c echo.Context) error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "func", "Suspend")
	appName := routines.GetProcessName(t.ClientName)
	err := routines.SuspendProcessWithCtx(ctx, appName)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("SuspendApp")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}

func (t *RoutineRequest) Resume(c echo.Context) error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "func", "Resume")
	appName := routines.GetProcessName(t.ClientName)
	err := routines.ResumeProcessWithCtx(ctx, appName)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("Resume: ResumeProcessWithCtx")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}
