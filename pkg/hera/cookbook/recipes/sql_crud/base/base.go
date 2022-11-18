package base

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/cookbook/recipes/common/sql_query/common"
	"github.com/zeus-fyi/olympus/pkg/hera/lib"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/vars"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

type ModelTemplate struct {
	lib.CodeGen
	*common.QueryMetadata
}

// NewPGModelTemplate should TODO make the parent pkg name part more seamless
// Should use this to create all Model templates to avoid any nil panics on query metadata
func NewPGModelTemplate(p filepaths.Path, queryInfo *common.QueryMetadata, dsnString string) ModelTemplate {
	if queryInfo == nil {
		queryInfo = &common.QueryMetadata{Type: "base", ParentPkgName: "models"}
	}
	queryInfo.ParentPkgName = "models"
	m := ModelTemplate{lib.NewCodeGen(p), queryInfo}
	m.NewInitPgConnToSchemaAutogen(dsnString)
	return m
}

// GenBaseImportHeaderLog to be used by structs which inherit this, shouldn't use for base logs without handling non-import pkgname prefix
func (m *ModelTemplate) GenBaseImportHeaderLog() *jen.Statement {
	modelVar := vars.NewVarGen()
	// modelVar.StringConstants["pkgName"] = "VarName" exported capital
	importedExternalVarName := "Sn"
	return common.GenPGDebugLogHeader(*m.QueryMetadata, modelVar.SetAndReturnImportedVarReference(m.ParentPkgName, importedExternalVarName))
}
