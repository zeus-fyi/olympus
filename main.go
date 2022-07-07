package main

import (
	"bitbucket.org/zeus/eth-indexer/cmd"
)

func main() {
	if err := cmd.Api(); err != nil {
		panic(err)
	}
}
