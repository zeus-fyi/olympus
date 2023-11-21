package hera_v1_codegen

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"

	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
)

var (
	sf = &strings_filter.FilterOpts{
		DoesNotStartWithThese: []string{"configs", "sandbox", "apps/external", ".git", ".circleci", ".DS_Store", ".idea", "apps/zeus/test/configs", "pkg/.DS_Store"},
		DoesNotInclude:        []string{"hardhat/artifacts", "node_modules", ".kube", "bin", "build", ".git", "hardhat/cache"},
	}
)

func ExtractSourceCode(ctx context.Context, f filepaths.Path) (CodeDirectoryMetadata, error) {
	// Step one: read the entire codebase
	cmd := CodeDirectoryMetadata{Map: make(map[string]CodeFilesMetadata)}
	// Function to recursively walk through directories
	var walkDir func(dir string, root string) error
	walkDir = func(dir string, root string) error {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			path := filepath.Join(dir, entry.Name())
			relPath, ferr := filepath.Rel(root, path)
			if ferr != nil {
				return ferr
			}
			if entry.IsDir() {
				ferr1 := walkDir(path, root)
				if ferr1 != nil {
					return ferr1
				}
			} else {
				// Read the contents of the file
				if strings_filter.FilterStringWithOpts(relPath, f.FilterFiles) {
					content, readErr := os.ReadFile(path)
					if readErr != nil {
						return readErr
					}
					cmd.SetContents(filepath.Dir(filepath.Clean(relPath)), relPath, content)
				}
			}
		}
		return nil
	}

	// Start walking from the root directory
	err := walkDir(f.DirIn, f.DirIn)
	if err != nil {
		panic(err)
	}
	aggregateDirectoryImports(&cmd)
	return cmd, err
}

func aggregateDirectoryImports(cmd *CodeDirectoryMetadata) {
	for dir, metadata := range cmd.Map {
		importSet := make(map[string]bool)
		for _, file := range metadata.GoCodeFiles.Files {
			for _, imp := range file.Imports {
				importSet[imp] = true
			}
		}
		var aggregatedImports []string
		for imp := range importSet {
			aggregatedImports = append(aggregatedImports, imp)
		}
		sort.Strings(aggregatedImports)
		// Update the entire struct in the map, not just the DirectoryImports field
		metadata.GoCodeFiles.DirectoryImports = aggregatedImports
		cmd.Map[dir] = metadata
	}
}

func (cm *CodeDirectoryMetadata) SetContents(dirIn, fn string, contents []byte) {
	cmdd, exists := cm.Map[dirIn]
	if !exists {
		cmdd = CodeFilesMetadata{}
	}
	baseFileName := filepath.Base(fn)
	switch {
	case strings.HasSuffix(fn, ".go"):
		goFileInfo, err := extractGoFileInfo(contents)
		if err != nil {
			panic(err)
		}
		if goFileInfo == nil {
			return
		}
		goFileInfo.FileName = baseFileName
		cmdd.GoCodeFiles.Files = append(cmdd.GoCodeFiles.Files, *goFileInfo)
	case strings.HasSuffix(fn, ".sql"):
		cmdd.SQLCodeFiles.Files = append(cmdd.SQLCodeFiles.Files, SQLCodeFile{
			FileName: baseFileName,
			Contents: string(contents),
		})
	case strings.HasSuffix(fn, ".yaml") || strings.HasSuffix(fn, ".yml"):
		cmdd.YamlCodeFiles.Files = append(cmdd.YamlCodeFiles.Files, YamlCodeFile{
			FileName: baseFileName,
			Contents: string(contents),
		})
	case strings.HasSuffix(fn, ".css"):
		cmdd.CssCodeFiles.Files = append(cmdd.CssCodeFiles.Files, CssCodeFile{
			FileName: baseFileName,
		})
	case strings.HasSuffix(fn, ".html"):
		cmdd.HtmlCodeFiles.Files = append(cmdd.HtmlCodeFiles.Files, HtmlCodeFile{
			FileName: baseFileName,
		})
	case strings.HasSuffix(fn, ".js") || strings.HasSuffix(fn, ".tsx"), strings.HasSuffix(fn, ".ts"):
		ext := ".js"
		if strings.HasSuffix(fn, ".tsx") {
			ext = ".tsx"
		}
		if strings.HasSuffix(fn, ".ts") {
			ext = ".ts"
		}
		cmdd.JsCodeFiles.Files = append(cmdd.JsCodeFiles.Files, JsCodeFile{
			FileName:  baseFileName,
			Contents:  string(contents),
			Extension: ext,
		})
	default:
		return
	}
	cm.Map[dirIn] = cmdd
}
