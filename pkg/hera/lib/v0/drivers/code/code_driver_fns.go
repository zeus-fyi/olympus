package code_driver

import (
	"github.com/zeus-fyi/jennifer/jen"
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

func (c *CodeDriverLib) AddSlice(jenCodeSlice []jen.Code) {
	if c.JenFile == nil {
		c.NewCodeGenFileShell()
	}
	for _, jc := range jenCodeSlice {
		c.JenFile.Add(jc)
	}
}
