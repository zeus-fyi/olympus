package main

import (
	"log"

	"github.com/zeus-fyi/olympus/cmd"
)

func main() {
	if err := cmd.ApiCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
