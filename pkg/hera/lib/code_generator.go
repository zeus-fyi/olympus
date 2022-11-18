package lib

import (
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/datastores"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives"
	code_driver "github.com/zeus-fyi/olympus/pkg/hera/lib/v0/drivers/code"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

type CodeGen struct {
	code_driver.CodeDriverLib
	primitives.PrimitiveGenerator
	datastores.DatastoreAutogen
}

func NewCodeGen(codeGenPath filepaths.Path) CodeGen {
	c := CodeGen{code_driver.NewCodeDriverLib(codeGenPath), primitives.PrimitiveGenerator{}, datastores.NewDatastoreAutogen()}
	return c
}
