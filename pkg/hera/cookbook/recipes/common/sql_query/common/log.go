package common

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/vars"
)

// GenPGDebugLogHeader crud -> to debug stamp.
// Use modelVar.GenImportedVarReference(importedPkgName) to get varName
func GenPGDebugLogHeader(metadata QueryMetadata, modelVar vars.VariableGen) *jen.Statement {
	var queryType string
	switch metadata.Type {
	case "read":
		queryType = "SelectQueryQuery"
	case "create":
		queryType = "InsertQuery"
	case "update":
		queryType = "UpdateQuery"
	case "delete":
		queryType = "DeleteQuery"
	default:
		queryType = "UndefinedQuery"
	}
	return jen.Id("log").Dot("Debug").Call().Dot("Interface").Call(jen.Lit(queryType+":"),
		jen.Id("q").Dot("LogHeader").Call(modelVar.GenImportedVarReference(metadata.ParentPkgName)))
}
