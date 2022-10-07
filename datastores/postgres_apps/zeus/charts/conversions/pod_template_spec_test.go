package conversions

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	v1 "k8s.io/api/apps/v1"
	v1core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TemplateSpecTestSuite struct {
	ChartPackagesTestSuite
}

func (s *TemplateSpecTestSuite) TestPodTemplateSpecConfigToDB1() {
	filepath := "/Users/alex/Desktop/Zeus/olympus/pkg/zeus/core/transformations/deployment.yaml"
	jsonBytes, err := yr.ReadYamlConfig(filepath)

	var d v1.Deployment
	err = json.Unmarshal(jsonBytes, &d)
	s.Require().Nil(err)

	podTemplateSpec := d.Spec.Template
	err = PodTemplateSpecConfigToDB(&podTemplateSpec)
	s.Require().Nil(err)
}

func (s *TemplateSpecTestSuite) TestPodTemplateSpecConfigToDB() {

	port := v1core.ContainerPort{
		Name:          "",
		HostPort:      0,
		ContainerPort: 0,
		Protocol:      "",
		HostIP:        "",
	}

	pc := v1core.Container{
		Name:            "",
		Image:           "",
		Command:         nil,
		Args:            nil,
		WorkingDir:      "",
		Ports:           []v1core.ContainerPort{port},
		EnvFrom:         nil,
		Env:             nil,
		Resources:       v1core.ResourceRequirements{},
		VolumeMounts:    nil,
		VolumeDevices:   nil,
		LivenessProbe:   nil,
		ReadinessProbe:  nil,
		StartupProbe:    nil,
		Lifecycle:       nil,
		SecurityContext: nil,
		Stdin:           false,
		StdinOnce:       false,
		TTY:             false,
	}

	podSpec := v1core.PodSpec{
		Volumes:          nil,
		InitContainers:   nil,
		Containers:       []v1core.Container{pc},
		SecurityContext:  nil,
		ImagePullSecrets: nil,
	}

	pst := v1core.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       podSpec,
	}

	err := PodTemplateSpecConfigToDB(&pst)
	s.Require().Nil(err)
}

func TestTemplateSpecTestSuite(t *testing.T) {
	suite.Run(t, new(TemplateSpecTestSuite))
}
