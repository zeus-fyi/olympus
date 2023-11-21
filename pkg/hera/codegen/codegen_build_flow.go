package hera_v1_codegen

import (
	"context"
	"fmt"

	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

const (
	DbSchemaDir   = "datastores/postgres/local_docker/docker-entrypoint-initdb.d/init.sql"
	PkgDir        = "pkg"
	AppsDir       = "apps"
	CookbooksDir  = "cookbooks"
	DockerDir     = "docker"
	DatastoresDir = "datastores"
	AiZeusApi     = "apps/olympus/zeus/api/v1/zeus/ai"
)

type BuildAiInstructions struct {
	GoFileDirs []string
}

func BuildAiInstructionsFromSourceCode(ctx context.Context, f filepaths.Path, buildDirs BuildAiInstructions) {
	sc, err := ExtractSourceCode(ctx, f)
	if err != nil {
		panic(err)
	}
	for _, bd := range buildDirs.GoFileDirs {
		for _, gf := range sc.Map[bd].GoCodeFiles.Files {
			fmt.Println(gf.PackageName, " | ", gf.FileName)
		}
	}
}
