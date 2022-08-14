package _func

import (
	"context"
	"fmt"
)

func funcName1(ctx context.Context, stringParam string) error {
	fmt.Println("Hello, world")
	return nil
}
