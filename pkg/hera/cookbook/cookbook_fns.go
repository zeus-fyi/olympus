package cookbook

import (
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

func (a *Cookbook) CreateTemplatesInPath(templatePath structs.Path) error {
	recipePath := a.ApplyCookBookToTemplatePath(templatePath)
	orderedTemplates := a.GetTopologicallySortedPaths(recipePath)
	for _, r := range orderedTemplates.Slice {
		err := a.CreateTemplate(r, CookbookPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Cookbook) ApplyCookBookToTemplatePath(templatePath structs.Path) structs.Path {
	templatePath.LeftExtendDirInPath(CookbookPath.DirIn)
	return templatePath
}

func (a *Cookbook) CustomZeusParsing(templatePath structs.Path) error {
	recipePath := a.ApplyCookBookToTemplatePath(templatePath)
	orderedTemplates := a.GetTopologicallySortedPaths(recipePath)
	for _, r := range orderedTemplates.Slice {
		err := a.CreateCustomZeusTemplate(r, CookbookPath)
		if err != nil {
			return err
		}
	}
	return nil
}
