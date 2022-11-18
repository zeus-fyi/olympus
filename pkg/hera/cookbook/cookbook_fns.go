package cookbook

import (
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

func (a *Cookbook) CreateTemplatesInPath(templatePath filepaths.Path) error {
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

func (a *Cookbook) ApplyCookBookToTemplatePath(templatePath filepaths.Path) filepaths.Path {
	templatePath.LeftExtendDirInPath(CookbookPath.DirIn)
	return templatePath
}

func (a *Cookbook) CustomZeusParsing(templatePath filepaths.Path) error {
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
