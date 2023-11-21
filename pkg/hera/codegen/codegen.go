package hera_v1_codegen

import (
	"bufio"
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
)

func CreateWorkflow(ctx context.Context, f filepaths.Path) (map[string]string, error) {
	// Step one: read the entire codebase
	fileMap := make(map[string]string)
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
					if strings.HasSuffix(relPath, ".DS_Store") {
						continue
					}
					content, readErr := ioutil.ReadFile(path)
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

	// Convert the map to JSON
	_, err = json.MarshalIndent(fileMap, "", "  ")
	if err != nil {
		panic(err)
	}
	return fileMap, err
}

type CodeDirectoryMetadata struct {
	Map map[string]CodeFilesMetadata
}

type CodeFilesMetadata struct {
	DirectoryImports []string
	GoCodeFiles      []GoCodeFile
}

type GoCodeFiles struct {
	Files            []GoCodeFile
	DirectoryImports []string
}

type GoCodeFile struct {
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
		for _, file := range metadata.GoCodeFiles {
			for _, imp := range file.Imports {
				importSet[imp] = true
			}
		}

		var aggregatedImports []string
		for imp := range importSet {
			aggregatedImports = append(aggregatedImports, imp)
		}
		// Update the entire struct in the map, not just the DirectoryImports field
		metadata.DirectoryImports = aggregatedImports
		cmd.Map[dir] = metadata
	}
}
func (cm *CodeDirectoryMetadata) SetContents(dirIn, fn string, contents []byte) {
	cmdd, exists := cm.Map[dirIn]
	if !exists {
		cmdd = CodeFilesMetadata{}
	}

	if strings.HasSuffix(fn, ".go") {
		packageName, imports := extractGoFileInfo(contents)
		cmdd.GoCodeFiles = append(cmdd.GoCodeFiles, GoCodeFile{
			PackageName: packageName,
			Imports:     imports,
			Contents:    string(contents),
		})
	}
	cm.Map[dirIn] = cmdd
}
