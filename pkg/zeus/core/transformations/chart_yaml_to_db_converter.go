package transformations

import (
	"fmt"
	"os"

	"github.com/ghodss/yaml"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/paths"
)

type YamlFileIO struct {
	chart_workload.TopologyBaseInfraWorkload
}

func (y *YamlFileIO) ReadK8sWorkloadDir(p filepaths.Path) error {
	err := paths.WalkAndApplyFuncToFileType(p, ".yaml", y.DecodeK8sWorkload)
	if err != nil {
		return err
	}
	return err
}

func (y *YamlFileIO) ReadYamlConfig(filepath string) ([]byte, error) {
	// Open YAML file
	jsonByteArray, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	jsonBytes, err := yaml.YAMLToJSON(jsonByteArray)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return jsonBytes, err
	}
	return jsonBytes, err
}

func (y *YamlFileIO) ReadK8sWorkloadInMemFsDir(p filepaths.Path, fs memfs.MemFS) error {
	err := fs.WalkAndApplyFuncToFileType(&p, ".yaml", y.DecodeK8sWorkloadFromInMemFS)
	if err != nil {
		return err
	}
	return err
}

func (y *YamlFileIO) DecodeK8sWorkload(filepath string) error {
	b, err := y.ReadYamlConfig(filepath)
	if err != nil {
		return err
	}
	err = y.DecodeBytes(b)
	return err
}

func (y *YamlFileIO) ReadYamlConfigInMemFS(filepath string, fs *memfs.MemFS) ([]byte, error) {
	// Open YAML file
	jsonByteArray, err := fs.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	jsonBytes, err := yaml.YAMLToJSON(jsonByteArray)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return jsonBytes, err
	}
	return jsonBytes, err
}

func (y *YamlFileIO) DecodeK8sWorkloadFromInMemFS(filepath string, fs *memfs.MemFS) error {
	b, err := y.ReadYamlConfigInMemFS(filepath, fs)
	if err != nil {
		return err
	}
	err = y.DecodeBytes(b)
	return err
}
