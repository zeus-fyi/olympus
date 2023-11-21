package hera_v1_codegen

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
)

func CreateWorkflow(ctx context.Context, f filepaths.Path) ([]byte, error) {
	// Step one: read the entire codebase
	fileMap := make(map[string]string)

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
			if entry.IsDir() && strings_filter.FilterStringWithOpts(relPath, f.FilterFiles) {
				ferr1 := walkDir(path, root)
				if ferr1 != nil {
					return ferr1
				}
			} else {
				fileMap[relPath] = path
			}
		}
		return nil
	}

	// Start walking from the root directory
	err := walkDir(f.DirIn, f.DirIn)
	if err != nil {
		panic(err)
	}

	// Convert the map to JSON
	jsonData, err := json.MarshalIndent(fileMap, "", "  ")
	if err != nil {
		panic(err)
	}
	return jsonData, err
}
