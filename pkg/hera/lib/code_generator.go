package lib

import (
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/datastores"
	code_driver "github.com/zeus-fyi/olympus/pkg/hera/lib/v0/drivers/code"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

type CodeGen struct {
	code_driver.CodeDriverLib
	primitives.PrimitiveGenerator
	datastores.DatastoreAutogen
}

func NewCodeGen(codeGenPath structs.Path) CodeGen {
	c := CodeGen{code_driver.NewCodeDriverLib(codeGenPath), primitives.PrimitiveGenerator{}, datastores.NewDatastoreAutogen()}
	return c
}
