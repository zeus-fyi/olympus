package common_fields

import "github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/fields"

func ErrField() fields.Field {
	return fields.Field{
		Pkg:   "",
		Name:  "err",
		Type:  "error",
		Value: "",
	}
}
