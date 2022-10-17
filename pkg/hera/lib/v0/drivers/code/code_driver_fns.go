package code_driver

import (
	"bytes"

	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io"
)

func (c *CodeDriverLib) NewCodeGenFileShell() {
	f := c.FileBaseGen.GenerateFileShell(c.Path)
	c.JenFile = f
}

func (c *CodeDriverLib) Add(jenCode jen.Code) {
	if c.JenFile == nil {
		c.NewCodeGenFileShell()
	}
	c.JenFile.Add(jenCode)
}

func (c *CodeDriverLib) AppendJenStatement(jenStmt *jen.Statement) {
	tmp := c.JenStatementChain
	tmp = append(tmp, jenStmt)
	c.JenStatementChain = tmp
}

func (c *CodeDriverLib) PopAndChainJenStatements() *jen.Statement {
	tmp := jen.Statement{}
	stmtChain := &tmp
	for _, stmt := range c.JenStatementChain {
		stmtChain.Add(stmt)
	}
	c.JenStatementChain = []*jen.Statement{}
	return stmtChain
}

var fileIO = file_io.FileIO{}

func (c *CodeDriverLib) Save() error {
	if c.JenFile == nil {
		c.NewCodeGenFileShell()
	}
	buf := &bytes.Buffer{}
	if err := c.JenFile.Render(buf); err != nil {
		return err
	}
	return fileIO.CreateFile(c.Path, buf.Bytes())
}
