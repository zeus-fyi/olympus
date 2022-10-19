package create

import (
	"github.com/zeus-fyi/jennifer/jen"
	common_fields "github.com/zeus-fyi/olympus/pkg/hera/cookbook/recipes/common/fields"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/fields"
)

func genInsertFnFields() []fields.Field {
	return []fields.Field{common_fields.CtxField(), common_fields.QueryParams()}
}

func WrappedErrField() fields.Field {
	return fields.Field{
		Pkg:     "",
		Name:    "err",
		Type:    "error",
		Value:   "",
		FnField: insertWrappedErrReturn(),
	}
}

func insertWrappedErrReturn() *jen.Statement {
	return jen.Id("misc").Dot("ReturnIfErr").Call(jen.Id("err"), jen.Id("q").Dot("LogHeader").Call(jen.Id("models").Dot("Sn")))
}

func genInsertFnReturnFields() []fields.Field {
	return []fields.Field{WrappedErrField()}
}
