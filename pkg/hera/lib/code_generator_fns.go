package lib

import (
	"bytes"

	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io"
)

func (cg *CodeGen) NewCodeGenFileShell() {
	f := cg.FileBaseGen.GenerateFileShell(cg.Path)
	cg.JenFile = f
}

func (cg *CodeGen) Add(jenCode jen.Code) {
	if cg.JenFile == nil {
		cg.NewCodeGenFileShell()
	}
	cg.JenFile.Add(jenCode)
}

func (cg *CodeGen) AppendJenStatement(jenStmt *jen.Statement) {
	tmp := cg.JenStatementChain
	tmp = append(tmp, jenStmt)
	cg.JenStatementChain = tmp
}

func (cg *CodeGen) PopAndChainJenStatements() *jen.Statement {
	tmp := jen.Statement{}
	stmtChain := &tmp
	for _, stmt := range cg.JenStatementChain {
		stmtChain.Add(stmt)
	}
	cg.JenStatementChain = []*jen.Statement{}
	return stmtChain
}

var fileIO = file_io.FileIO{}

func (cg *CodeGen) Save() error {
	if cg.JenFile == nil {
		cg.NewCodeGenFileShell()
	}
	buf := &bytes.Buffer{}
	if err := cg.JenFile.Render(buf); err != nil {
		return err
	}
	return fileIO.CreateFile(cg.Path, buf.Bytes())
}
