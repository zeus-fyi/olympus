package code_driver

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/file_shells/base"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

type CodeDriverLib struct {
	Path              structs.Path
	FileBaseGen       base.FileComponentBaseElements
	JenStatementChain []*jen.Statement
	JenFile           *jen.File
}

func NewCodeDriverLib(codeGenPath structs.Path) CodeDriverLib {
	c := CodeDriverLib{
		Path:              codeGenPath,
		FileBaseGen:       base.FileComponentBaseElements{},
		JenStatementChain: []*jen.Statement{},
		JenFile:           nil,
	}
	return c
}
