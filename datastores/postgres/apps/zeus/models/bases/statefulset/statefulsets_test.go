package statefulset

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

type StatefulSetTestSuite struct {
	base.TestSuite
	TestDirectory string
}

func (s *StatefulSetTestSuite) SetupTest() {
	s.TestDirectory = "./statefulset.yaml"
}
func (s *StatefulSetTestSuite) TestStatefulSetK8sToDBConversion() {
	sts := NewStatefulSet()
	filepath := s.TestDirectory
	jsonBytes, err := ReadYamlConfig(filepath)
	s.Require().Nil(err)
	err = json.Unmarshal(jsonBytes, &sts.K8sStatefulSet)
	s.Require().Nil(err)
	s.Require().NotEmpty(sts.K8sStatefulSet)

	err = sts.ConvertK8sStatefulSetToDB()
	s.Require().Nil(err)
	s.Require().NotEmpty(sts.Spec)
	s.Require().NotEmpty(sts.Metadata)
	s.Assert().Equal("name", sts.Metadata.Name.ChartSubcomponentKeyName)
	s.Assert().Equal("zeus-lighthouse", sts.Metadata.Name.ChartSubcomponentValue)

	s.Require().NotEmpty(sts.Spec.Replicas)
	s.Assert().Equal("replicas", sts.Spec.Replicas.ChartSubcomponentKeyName)
	s.Assert().Equal("1", sts.Spec.Replicas.ChartSubcomponentValue)

	s.Require().NotEmpty(sts.Spec.Selector)

	s.Assert().Equal("selector", sts.Spec.Selector.MatchLabels.ChartSubcomponentChildClassTypeName)
	selectorValues := sts.Spec.Selector.MatchLabels.Values
	s.Assert().Len(selectorValues, 1)

	for _, ml := range selectorValues {
		s.Assert().Equal("selectorString", ml.ChartSubcomponentKeyName)
		expectedSelectorMatchLabels := "{\"matchLabels\":{\"app.kubernetes.io/instance\":\"zeus\",\"app.kubernetes.io/name\":\"lighthouse\"}}"
		s.Assert().Equal(expectedSelectorMatchLabels, ml.ChartSubcomponentValue)
	}
	s.Assert().Equal("zeus-lighthouse", sts.Metadata.Name.ChartSubcomponentValue)

	s.Require().NotEmpty(sts.Spec.PodManagementPolicy)
	s.Assert().Equal("podManagementPolicy", sts.Spec.PodManagementPolicy.ChartSubcomponentKeyName)
	s.Assert().Equal("OrderedReady", sts.Spec.PodManagementPolicy.ChartSubcomponentValue)

	s.Require().NotEmpty(sts.Spec.Template)
	s.Require().NotEmpty(sts.Spec.Template.Metadata)

	s.Assert().Equal("labels", sts.Spec.Template.Metadata.Labels.ChartSubcomponentChildClassTypeName)

	templateSpecMetadataLabelValues := sts.Spec.Template.Metadata.Labels.Values

	// zeus adds a version label
	s.Assert().Len(templateSpecMetadataLabelValues, 3)
	countLabels := 0
	for _, label := range sts.Spec.Template.Metadata.Labels.Values {
		if label.ChartSubcomponentKeyName == "version" && strings.HasPrefix(label.ChartSubcomponentValue, "version-") {
			countLabels += 1
		}
		if label.ChartSubcomponentKeyName == "app.kubernetes.io/name" && label.ChartSubcomponentValue == "lighthouse" {
			countLabels += 10
		}
		if label.ChartSubcomponentKeyName == "app.kubernetes.io/instance" && label.ChartSubcomponentValue == "zeus" {
			countLabels += 100
		}
	}
	s.Assert().Equal(111, countLabels)
	s.Require().NotEmpty(sts.Spec.VolumeClaimTemplates)
	s.Require().NotEmpty(sts.Spec.Template.Spec.PodTemplateContainers)

	count := 0
	for _, pvc := range sts.Spec.VolumeClaimTemplates.VolumeClaimTemplateSlice {
		s.Assert().Equal("storage", pvc.Metadata.Metadata.Name.ChartSubcomponentValue)

		expectKeyNameAccessMode := "accessMode"
		expectKeyNameRequests := "requests"

		s.Assert().Equal("storageClassName", pvc.Spec.StorageClassName.ChartSubcomponentKeyName)
		s.Assert().Equal("beaconStorageClassName", pvc.Spec.StorageClassName.ChartSubcomponentValue)
		for _, rr := range pvc.Spec.ResourceRequests.Values {
			if rr.ChartSubcomponentKeyName == expectKeyNameRequests {
				val := strings.Trim(rr.ChartSubcomponentValue, `""`)
				s.Assert().Equal("20Gi", val)
				count += 10
			}
		}
		for _, am := range pvc.Spec.AccessModes.Values {
			if am.ChartSubcomponentKeyName == expectKeyNameAccessMode {
				s.Assert().Equal("ReadWriteOnce", am.ChartSubcomponentValue)
				count += 1
			}
		}
	}

	s.Assert().Equal(11, count)
	c := charts.NewChart()
	ts := chronos.Chronos{}
	c.ChartPackageID = ts.UnixTimeStampNow()
	subCTEs := sts.GetStatefulSetCTE(&c)

	s.Assert().NotEmpty(subCTEs)

	fmt.Println(subCTEs.GenerateChainedCTE())
	fmt.Println(subCTEs.Params)
}

func TestStatefulSetTestSuite(t *testing.T) {
	suite.Run(t, new(StatefulSetTestSuite))
}

func ReadYamlConfig(filepath string) ([]byte, error) {
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
