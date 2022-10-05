package zeus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

func ConvertYamlConfig(filepath string) error {
	// Open YAML file
	jsonByteArray, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	jsonBytes, err := yaml.YAMLToJSON(jsonByteArray)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	m := make(map[string]interface{})

	err = json.Unmarshal(jsonBytes, &m)
	return err
}

// should match this format on query
//{
//	"apiVersion":"v1",
//  "kind":"Service",
//	"metadata":{
//				"labels":{"app.kubernetes.io/instance":"s"},
//				"name":"s"
//				},
//	"spec": {
//				"ports":[{
//							"name":"http",
//							"port":80,
//							"protocol":"TCP",
//							"targetPort":"http"
//						}],
//				"selector": {"app.kubernetes.io/instance":"s"},
//				"type":"ClusterIP"
//			 }
//}
