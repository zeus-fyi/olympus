package primitives

import (
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/conditionals"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/funcs"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/structs"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/vars"
)

type PrimitiveGenerator struct {
	Fns       funcs.FuncGen
	Vars      vars.VariableGen
	Structs   structs.StructsGen
	CaseStmts conditionals.SwitchCase
}
