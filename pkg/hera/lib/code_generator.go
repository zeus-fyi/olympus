package lib

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/file_shells/base"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

type CodeGen struct {
	Fw          primitives.FileWrapper
	Path        structs.Path
	FileBaseGen base.FileComponentBaseElements

	StructsToGen []primitives.StructGen
	FuncToGen    []primitives.FuncGen

	JenStatementChain []*jen.Statement
	JenFile           *jen.File
}

func NewCodeGen(codeGenPath structs.Path) CodeGen {
	return CodeGen{
		FileBaseGen:       base.FileComponentBaseElements{},
		Path:              codeGenPath,
		StructsToGen:      []primitives.StructGen{},
		FuncToGen:         []primitives.FuncGen{},
		JenStatementChain: []*jen.Statement{},
	}
}
