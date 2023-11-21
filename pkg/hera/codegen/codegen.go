package hera_v1_codegen

import (
	"bytes"
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
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
	Functions   map[string]FunctionInfo
	Contents    string
}

type FunctionInfo struct {
	Name       string
	Parameters string
	ReturnType string
	Body       string
}

func extractGoFileInfoV2(src interface{}) (*GoCodeFile, error) {
	fset := token.NewFileSet() // positions are relative to fset
	fileInfo := &GoCodeFile{
		Functions: make(map[string]FunctionInfo),
	}
	var reader *strings.Reader
	switch s := src.(type) {
	case string:
		reader = strings.NewReader(s)
		fileInfo.Contents = s
	case []byte:
		reader = strings.NewReader(string(s))
		fileInfo.Contents = string(s)
	default:
		return nil, fmt.Errorf("src must be a string or []byte")
	}
	// Parse the file given in arguments
	// Parse the source code
	node, err := parser.ParseFile(fset, "", reader, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// Package name
	fileInfo.PackageName = node.Name.Name

	// Inspect the AST and find all imports and functions
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ImportSpec:
			if x.Path != nil {
				fileInfo.Imports = append(fileInfo.Imports, x.Path.Value)
			}
		case *ast.FuncDecl:
			var funcInfo FunctionInfo
			if x.Name != nil {
				funcInfo.Name = x.Name.Name
			}

			// Parameters and Return Type
			funcInfo.Parameters = fieldListToString(fset, x.Type.Params)
			funcInfo.ReturnType = fieldListToString(fset, x.Type.Results)

			// Function body
			var body bytes.Buffer
			if x.Body != nil {
				err = printer.Fprint(&body, fset, x.Body)
				if err != nil {
					panic(err)
				}
				funcInfo.Body = body.String()
			}
			fileInfo.Functions[funcInfo.Name] = funcInfo
		}
		return true
	})

	return fileInfo, nil
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
		goFileInfo, err := extractGoFileInfoV2(contents)
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

func fieldListToString(fset *token.FileSet, fl *ast.FieldList) string {
	if fl == nil {
		return ""
	}

	var buf bytes.Buffer
	for _, field := range fl.List {
		if len(field.Names) > 0 {
			for _, name := range field.Names {
				buf.WriteString(name.Name)
				buf.WriteString(" ")
			}
		}

		err := printer.Fprint(&buf, fset, field.Type)
		if err != nil {
			panic(err)
		}
		buf.WriteString(", ")
	}

	// Remove trailing comma and space
	result := buf.String()
	if len(result) >= 2 {
		result = result[:len(result)-2]
	}
	return result
}
