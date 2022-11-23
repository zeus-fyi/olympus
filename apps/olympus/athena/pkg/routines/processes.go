package routines

import (
	"context"
	"fmt"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v3/process"
)

func KillProcessWithCtx(ctx context.Context, name string) error {
	processes, err := process.Processes()
	if err != nil {
		return err
	}
	for _, p := range processes {
		n, perr := p.Name()
		if perr != nil {
			log.Err(err).Msgf("killProcessWithCtx %s", name)
			return err
		}
		if n == name {
			return p.SendSignalWithContext(ctx, syscall.SIGINT)
		}
	}
	return fmt.Errorf("process not found")
}

func SuspendProcessWithCtx(ctx context.Context, name string) error {
	processes, err := process.Processes()
	if err != nil {
		return err
	}
	for _, p := range processes {
		n, perr := p.Name()
		if perr != nil {
			log.Err(err).Msgf("suspendProcessWithCtx %s", name)
			return err
		}
		if n == name {
			return p.SuspendWithContext(ctx)
		}
	}
	return fmt.Errorf("process not found")
}
