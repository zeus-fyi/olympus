package _func

import (
	"context"
	"errors"
)

func templateFunc(ctx context.Context, param string) error {
	if len(param) <= 0 {
		return errors.New("error message")
	}
	return nil
}
