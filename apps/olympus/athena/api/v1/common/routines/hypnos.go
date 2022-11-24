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

func (t *RoutineRequest) HypnosSuspend(c echo.Context) error {
	appName := routines.GetProcessName(t.ClientName)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "func", "HypnosSuspend")
	err := routines.SuspendProcessWithCtx(ctx, appName)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("HypnosSuspend: SuspendProcessWithCtx")
		return c.JSON(http.StatusInternalServerError, err)
	}
	err = t.InjectHypnos()
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("HypnosSuspend: InjectHypnos")
		return c.JSON(http.StatusInternalServerError, err)
	}
	resp := RoutineResp{Status: "hypnos started"}
	return c.JSON(http.StatusOK, resp)
}

func (t *RoutineRequest) InjectHypnos() error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "func", "InjectHypnos")
	port := routines.GetProcessPorts(t.ClientName)
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("hypnos --port=%s", port))
	err := cmd.Start()
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("InjectHypnos")
		return err
	}
	return nil
}
