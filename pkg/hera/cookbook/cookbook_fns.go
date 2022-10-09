package cookbook

import (
	"github.com/zeus-fyi/olympus/pkg/utils/printer/v0/structs"
	"golang.org/x/exp/maps"
)

func (a *Cookbook) ApplyCookBookToTemplatePath(templatePath structs.Path) structs.Path {
	templatePath.LeftExtendDirInPath(CookbookPath.DirIn)
	return templatePath
}

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

func (a *Cookbook) GetTopologicallySortedPaths(path structs.Path) structs.Paths {
	templatePathMaps := print.BuildPathsFromDirInPath(path, ".go")
	depth := len(maps.Keys(templatePathMaps))
	tmp := structs.Paths{}
	for i := 0; i <= depth; i++ {
		recipePaths := templatePathMaps[i]
		for _, r := range recipePaths.Slice {
			tmp.AddPathToSlice(r)
		}
		i += i
	}
	return tmp
}
