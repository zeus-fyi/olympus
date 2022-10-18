package funcs

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/structs"
)

func (f *FuncGen) GenerateStructPtrFunc(s structs.StructGen) jen.Code {
	if len(s.Name) <= 0 {
		return jen.Nil()
	}
	header := jen.Func().Params(jen.Id(s.ShortHand()).Op("*").Id(s.Name))
	fn := f.GenerateFuncShell(header)
	return fn
}
