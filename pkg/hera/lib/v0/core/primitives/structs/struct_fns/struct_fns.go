package struct_fns

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/funcs"
)

type StructFn struct {
	funcs.FuncGen
}

func (s *StructFn) GenerateStructPtrFunc(structName string) jen.Code {
	if len(structName) <= 0 {
		return jen.Nil()
	}
	fn := s.GenerateStructFunc(structName)
	return fn
}
