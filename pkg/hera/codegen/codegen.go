package hera_v1_codegen

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"

	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
)

const (
	DbSchemaDir   = "datastores/postgres/local_docker/docker-entrypoint-initdb.d/init.sql"
	PkgDir        = "pkg"
	AppsDir       = "apps"
	CookbooksDir  = "cookbooks"
	DockerDir     = "docker"
	DatastoresDir = "datastores"
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

type CodeDirectoryMetadata struct {
	Map map[string]CodeFilesMetadata
}

type CodeFilesMetadata struct {
	GoCodeFiles   GoCodeFiles
	SQLCodeFiles  SQLCodeFiles
	YamlCodeFiles YamlCodeFiles
	JsCodeFiles   JsCodeFiles
	CssCodeFiles  CssCodeFiles
	HtmlCodeFiles HtmlCodeFiles
}

type JsCodeFiles struct {
	Files            []JsCodeFile
	DirectoryImports []string
}

type JsCodeFile struct {
	FileName  string
	Extension string
	Imports   []string
	Contents  string
}
type CssCodeFiles struct {
	Files []CssCodeFile
}

type CssCodeFile struct {
	FileName string
}
type HtmlCodeFiles struct {
	Files []HtmlCodeFile
}

type HtmlCodeFile struct {
	FileName string
}

type SQLCodeFiles struct {
	Files []SQLCodeFile
}

type SQLCodeFile struct {
	FileName string
	Contents string
}

type YamlCodeFiles struct {
	Files []YamlCodeFile
}

type YamlCodeFile struct {
	FileName string
	Contents string
}

type GoCodeFiles struct {
	Files            []GoCodeFile
	DirectoryImports []string
}

type GoCodeFile struct {
	FileName    string
	PackageName string
	Imports     []string
	Contents    string
}

func extractGoFileInfo(contents []byte) (string, []string) {
	scanner := bufio.NewScanner(strings.NewReader(string(contents)))
	var packageName string
	var imports []string
	var inImportBlock bool

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "package ") {
			packageName = strings.TrimSpace(strings.TrimPrefix(line, "package"))
		} else if strings.HasPrefix(line, "import (") {
			inImportBlock = true
		} else if inImportBlock && line == ")" {
			inImportBlock = false
		} else if inImportBlock {
			// Remove the quotes around the import path
			trimmedLine := strings.Trim(line, ` "`)
			if trimmedLine != "" {
				imports = append(imports, trimmedLine)
			}
		} else if strings.HasPrefix(line, "import ") {
			// Single import
			importedPackage := strings.Trim(strings.TrimPrefix(line, "import "), `"`)
			imports = append(imports, importedPackage)
		}
	}

	if err := scanner.Err(); err != nil {
		// Handle the error according to your needs
	}
	return packageName, imports
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
		packageName, imports := extractGoFileInfo(contents)
		cmdd.GoCodeFiles.Files = append(cmdd.GoCodeFiles.Files, GoCodeFile{
			FileName:    baseFileName,
			PackageName: packageName,
			Imports:     imports,
			Contents:    string(contents),
		})
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
