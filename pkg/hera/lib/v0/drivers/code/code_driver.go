package code_driver

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/file_shells/base"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

type CodeDriverLib struct {
	Path              filepaths.Path
	FileBaseGen       base.FileComponentBaseElements
	JenStatementChain []*jen.Statement
	JenFile           *jen.File
}

func NewCodeDriverLib(codeGenPath filepaths.Path) CodeDriverLib {
	c := CodeDriverLib{
		Path:              codeGenPath,
		FileBaseGen:       base.FileComponentBaseElements{},
		JenStatementChain: []*jen.Statement{},
		JenFile:           nil,
	}
	return c
}

func (c *CodeDriverLib) ResetInternalJenCaches() {
	c.FileBaseGen = base.FileComponentBaseElements{}
	c.JenStatementChain = []*jen.Statement{}
	c.JenFile = nil
}
