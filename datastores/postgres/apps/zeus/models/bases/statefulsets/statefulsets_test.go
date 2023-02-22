package statefulsets

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type StatefulSetTestSuite struct {
	test_suites_base.TestSuite
	TestDirectory string
}

func (s *StatefulSetTestSuite) SetupTest() {
	s.TestDirectory = "./statefulsets.yaml"
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
	s.Require().NotEmpty(sts.Metadata)
	s.Assert().Equal("name", sts.Metadata.Name.ChartSubcomponentKeyName)
	s.Assert().Equal("zeus-consensus-client", sts.Metadata.Name.ChartSubcomponentValue)

	s.Require().NotEmpty(sts.Spec)

	// podManagementPolicy
	pmp := sts.Spec.PodManagementPolicy
	s.Require().NotEmpty(pmp)
	s.Assert().Equal("StatefulSetPodManagementPolicy", pmp.ChartSubcomponentChildClassTypeName)
	s.Assert().Equal("podManagementPolicy", pmp.ChartSubcomponentKeyName)
	s.Assert().Equal("OrderedReady", pmp.ChartSubcomponentValue)

	// replicas
	replicas := sts.Spec.Replicas
	s.Require().NotEmpty(replicas)
	s.Assert().Equal("replicas", replicas.ChartSubcomponentChildClassTypeName)
	s.Assert().Equal("replicas", replicas.ChartSubcomponentKeyName)
	s.Assert().Equal("1", replicas.ChartSubcomponentValue)

	// selector
	s.Require().NotEmpty(sts.Spec.Selector)
	s.Assert().Equal("selector", sts.Spec.Selector.MatchLabels.ChartSubcomponentChildClassTypeName)
	selectorValues := sts.Spec.Selector.MatchLabels.Values
	s.Assert().Len(selectorValues, 1)

	for _, ml := range selectorValues {
		s.Assert().Equal("selectorString", ml.ChartSubcomponentKeyName)
		expectedSelectorMatchLabels := "{\"matchLabels\":{\"app.kubernetes.io/instance\":\"zeus-consensus-client\",\"app.kubernetes.io/name\":\"zeus-consensus-client\"}}"
		s.Assert().Equal(expectedSelectorMatchLabels, ml.ChartSubcomponentValue)
	}
	s.Assert().Equal("zeus-consensus-client", sts.Metadata.Name.ChartSubcomponentValue)

	// serviceName
	svcName := sts.Spec.ServiceName
	s.Assert().NotEmpty(svcName)
	s.Assert().Equal("StatefulSetServiceName", svcName.ChartSubcomponentChildClassTypeName)
	s.Assert().Equal("serviceName", svcName.ChartSubcomponentKeyName)
	s.Assert().Equal("zeus-consensus-client", svcName.ChartSubcomponentValue)

	// strategy
	stsStrategy := sts.Spec.StatefulSetUpdateStrategy
	s.Assert().NotEmpty(stsStrategy)
	s.Assert().Equal("StatefulSetUpdateStrategy", stsStrategy.ChartSubcomponentChildClassTypeName)

	stsStrategyValues := stsStrategy.Values
	s.Assert().NotEmpty(stsStrategyValues)
	stratCount := 0
	for _, strat := range stsStrategyValues {
		if strat.ChartSubcomponentKeyName == "type" && strat.ChartSubcomponentValue == "RollingUpdate" {
			stratCount += 1
		}
	}
	s.Assert().Equal(1, stratCount)

	// pod template spec
	stsSpecTemplate := sts.Spec.Template
	s.Require().NotEmpty(stsSpecTemplate)

	s.Require().NotNil(stsSpecTemplate.Spec.ShareProcessNamespace)
	s.Assert().Equal("true", stsSpecTemplate.Spec.ShareProcessNamespace.ChartSubcomponentValue)

	// template spec labels
	s.Assert().Equal("labels", stsSpecTemplate.Metadata.Labels.ChartSubcomponentChildClassTypeName)
	templateSpecMetadataLabelValues := stsSpecTemplate.Metadata.Labels.Values

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
	// todo fix this below
	s.Assert().Equal(1, countLabels)

	// containers
	conts := sts.Spec.Template.Spec.PodTemplateContainers
	s.Assert().NotEmpty(conts)
	s.Assert().Len(conts, 4)

	countVerifiedCont := 0
	countVerifiedVolMountCont := 0
	countVerifiedProbes := 0
	countInitContainers := 0
	for _, cont := range conts {
		name := cont.Metadata.ContainerName
		imageID := cont.Metadata.ContainerImageID
		imagePullPolicy := cont.Metadata.ContainerImagePullPolicy

		cmdToRun := cont.CmdArgs.CommandValues
		cmdArgs := cont.CmdArgs.ArgsValues
		contVolMounts := cont.VolumeMounts
		contSecCtx := cont.SecurityContext
		contEnvs := cont.Env

		contProbes := cont.Probes
		if cont.Metadata.IsInitContainer && name == "init-chown-data" {
			countInitContainers += 1
			s.Assert().Equal("init-chown-data", name)
			s.Assert().Equal("busybox:1.34.0", imageID)
			s.Assert().Equal("chown,-R,10001:10001,/data", cmdToRun)
			s.Assert().Equal("IfNotPresent", imagePullPolicy)
			s.Assert().NotEmpty(contSecCtx)
			s.Assert().Equal("{\"runAsUser\":0,\"runAsNonRoot\":false}", contSecCtx.SecurityContextKeyValues)

			s.Assert().NotEmpty(contVolMounts)
			for _, cvm := range contVolMounts {
				s.Assert().Equal("consensus-client-storage", cvm.VolumeName)
				s.Assert().Equal("/data", cvm.VolumeMountPath)
				countVerifiedVolMountCont += 1
			}
			countVerifiedCont += 1
		}

		if cont.Metadata.IsInitContainer && name == "init-snapshot" {
			countInitContainers += 1
			countVerifiedCont += 1
			s.Assert().Len(cont.Env, 1)
		}

		if cont.Metadata.IsInitContainer == false && name == "lighthouse" {
			s.Assert().Equal("lighthouse", name)
			s.Assert().Equal("sigp/lighthouse:v3.1.0", imageID)
			s.Assert().Equal("/bin/sh", cmdToRun)
			s.Assert().Equal("-c,/scripts/start.sh", cmdArgs)
			s.Assert().Equal("IfNotPresent", imagePullPolicy)

			s.Assert().NotEmpty(cont.Ports)
			s.Assert().NotEmpty(contProbes)

			for _, pr := range contProbes {

				if pr.ProbeType == "livenessProbe" {
					s.Assert().Equal("{\"tcpSocket\":{\"port\":\"http-api\"},\"initialDelaySeconds\":60,\"periodSeconds\":120}", pr.ProbeKeyValuesJSONb)
					countVerifiedProbes += 1
				}
				if pr.ProbeType == "readinessProbe" {
					s.Assert().Equal("{\"tcpSocket\":{\"port\":\"http-api\"},\"initialDelaySeconds\":10,\"periodSeconds\":10}", pr.ProbeKeyValuesJSONb)
					countVerifiedProbes += 10
				}
			}
			s.Assert().NotEmpty(contEnvs)

			for _, env := range contEnvs {
				s.Assert().Equal("POD_IP", env.Name)
				s.Assert().Equal("{\"fieldRef\":{\"fieldPath\":\"status.podIP\"}}", env.Value)
			}

			s.Require().NotEmpty(contVolMounts)
			for _, cvm := range contVolMounts {
				cvmName := cvm.VolumeName
				if cvmName == "storage" {
					s.Assert().Equal("/data", cvm.VolumeMountPath)
					countVerifiedVolMountCont += 10
				}
				if cvmName == "jwt" {
					s.Assert().Equal("/data/jwt.hex", cvm.VolumeMountPath)
					s.Assert().Equal("jwt.hex", cvm.VolumeSubPath)
					s.Assert().True(cvm.VolumeReadOnly)
					countVerifiedVolMountCont += 100
				}
			}
			countVerifiedCont += 10
		}
	}

	s.Assert().Equal(2, countInitContainers)
	// todo fix
	//s.Assert().Equal(11, countVerifiedProbes)
	//s.Assert().Equal(11, countVerifiedCont)
	//s.Assert().Equal(111, countVerifiedVolMountCont)

	// pod spec volumes
	podSpecVolumes := sts.Spec.Template.Spec.PodTemplateSpecVolumes
	s.Assert().NotEmpty(podSpecVolumes)
	//s.Assert().Len(podSpecVolumes, 2)

	//volCount := 0
	//for _, psv := range podSpecVolumes {
	//	if psv.VolumeName == "jwt" && psv.VolumeKeyValuesJSONb == "{\"name\":\"jwt\",\"secret\":{\"secretName\":\"zeus-lighthouse-jwt\"}}" {
	//		volCount += 1
	//	}
	//	if psv.VolumeName == "storage" && psv.VolumeKeyValuesJSONb == "{\"name\":\"storage\",\"emptyDir\":{}}" {
	//		volCount += 10
	//	}
	//}
	//s.Assert().Equal(11, volCount)

	s.Require().NotEmpty(sts.Spec.Template.Metadata)

	s.Require().NotEmpty(sts.Spec.VolumeClaimTemplates)
	s.Require().NotEmpty(sts.Spec.Template.Spec.PodTemplateContainers)

	//count := 0
	//for _, pvc := range sts.Spec.VolumeClaimTemplates.VolumeClaimTemplateSlice {
	//	s.Assert().Equal("storage", pvc.Metadata.Metadata.Name.ChartSubcomponentValue)
	//
	//	expectKeyNameAccessMode := "accessMode"
	//	expectKeyNameRequests := "requests"
	//
	//	s.Assert().Equal("storageClassName", pvc.Spec.StorageClassName.ChartSubcomponentKeyName)
	//	s.Assert().Equal("beaconStorageClassName", pvc.Spec.StorageClassName.ChartSubcomponentValue)
	//	for _, rr := range pvc.Spec.ResourceRequests.Values {
	//		if rr.ChartSubcomponentKeyName == expectKeyNameRequests {
	//			val := strings.Trim(rr.ChartSubcomponentValue, `""`)
	//			s.Assert().Equal("300Gi", val)
	//			count += 10
	//		}
	//	}
	//	for _, am := range pvc.Spec.AccessModes.Values {
	//		if am.ChartSubcomponentKeyName == expectKeyNameAccessMode {
	//			s.Assert().Equal("ReadWriteOnce", am.ChartSubcomponentValue)
	//			count += 1
	//		}
	//	}
	//}

	//s.Assert().Equal(11, count)
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
