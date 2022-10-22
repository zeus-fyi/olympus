package structs

import (
	"strings"

	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/fields"
)

type StructGen struct {
	Name       string
	Fields     []fields.Field
	PluralDecl jen.Code
}

func (s *StructGen) GenerateSliceType() jen.Code {
	s.PluralDecl = jen.Null().Type().Id(s.Name + "s").Index().Id(s.Name)
	return s.PluralDecl
}

func (s *StructGen) ShortHand() string {
	if len(s.Name) > 0 {
		return strings.ToLower(s.Name[0:1])
	}
	return ""
}
