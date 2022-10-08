package hera

import (
	"github.com/zeus-fyi/olympus/pkg/hera/lib/drivers/template"
	"github.com/zeus-fyi/olympus/pkg/utils/dev_hacks"
	"github.com/zeus-fyi/olympus/pkg/utils/printer"
	"github.com/zeus-fyi/olympus/pkg/utils/printer/v0/structs"
)

var (
	CookbookPath = structs.Path{
		DirIn:  "cookbook/code_templates",
		DirOut: "cookbook/autogen",
	}
	p = printer.Printer{}
)

func CreateTemplate(templatePath structs.Path) error {
	recipePath := ApplyCookBookToTemplatePath(templatePath)
	err := template.CreateJenFile(recipePath)
	return err
}

func ApplyCookBookToTemplatePath(templatePath structs.Path) structs.Path {
	templatePath.LeftExtendDirInPath(CookbookPath.DirIn)
	templatePath.LeftExtendDirOutPath(CookbookPath.DirOut)
	return templatePath
}
func CreateTemplatesInPath(templatePath structs.Path) {
	recipePath := ApplyCookBookToTemplatePath(templatePath)
	templatePaths := GetPaths(recipePath)

	dev_hacks.Use(templatePaths)
}

func GetPaths(path structs.Path) structs.Paths {
	templatePaths := p.BuildPathsFromDirInPath(path, ".go")

	return templatePaths
}
