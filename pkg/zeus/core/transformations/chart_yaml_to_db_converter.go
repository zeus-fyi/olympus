package transformations

import (
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/paths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

type YamlReader struct {
	chart_workload.NativeK8s
}

func (y *YamlReader) ReadK8sWorkloadDir(p structs.Path) error {
	err := paths.WalkAndApplyFuncToFileType(p.DirIn, ".yaml", y.DecodeK8sWorkload)
	if err != nil {
		return err
	}
	return err
}

func (y *YamlReader) ReadYamlConfig(filepath string) ([]byte, error) {
	// Open YAML file
	jsonByteArray, err := ioutil.ReadFile(filepath)
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

func (y *YamlReader) DecodeK8sWorkload(filepath string) error {
	b, err := y.ReadYamlConfig(filepath)
	if err != nil {
		return err
	}
	err = y.DecodeBytes(b)
	return err
}
