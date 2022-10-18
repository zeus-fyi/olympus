package common_fields

import "github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/fields"

func CtxField() fields.Field {
	return fields.Field{
		Pkg:   "context",
		Name:  "ctx",
		Type:  "Context",
		Value: "",
	}
}
