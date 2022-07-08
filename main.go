package main

import (
	"github.com/zeus-fyi/olympus/cmd"
)

func main() {
	if err := cmd.Api(); err != nil {
		panic(err)
	}
}
