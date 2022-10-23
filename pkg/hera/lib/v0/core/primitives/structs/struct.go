package structs

import (
	"strings"

	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/fields"
	"github.com/zeus-fyi/tables-to-go/pkg/database"
)

type StructGen struct {
	Name       string
	Fields     []fields.Field
	PluralDecl jen.Code

	DBTable *database.Table
}

func (s *StructGen) GenerateSliceType() jen.Code {
	s.PluralDecl = jen.Null().Type().Id(s.Name + "Slice").Index().Id(s.Name)
	return s.PluralDecl
}

func (s *StructGen) ShortHand() string {
	if len(s.Name) > 0 {
		return strings.ToLower(s.Name[0:1])
	}
	return ""
}

func (s *StructGen) TableExpressionName() string {
	if s.DBTable != nil && len(s.DBTable.Name) > 0 {
		return s.DBTable.Name
	}
	return ""
}
