package transformations

import (
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/writers"
)

func (y *YamlFileIO) WriteYamlConfig(p filepaths.Path, jsonByteArray []byte) error {
	jsonBytes, err := yaml.JSONToYAML(jsonByteArray)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	w := writers.WriterLib{}
	err = w.CreateV2FileOut(p, jsonBytes)
	return err
}
