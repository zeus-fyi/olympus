package containers

import (
	"encoding/json"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/containers"
	v1 "k8s.io/api/core/v1"
)

func (c *Container) ParseFields() error {
	env, err := c.DB.parseEnvVars(c.DB.EnvVar)
	if err != nil {
		return err
	}
	c.K8sContainer.Env = env
	portNum := int32(c.DB.PortNumber)
	ports := c.DB.parseContainerPorts(c.DB.PortName, c.DB.PortProtocol, portNum)
	c.K8sContainer.Ports = append(c.K8sContainer.Ports, ports...)

	vm, err := c.DB.parseVolumeMount(c.DB.VolumeName, c.DB.VolumePath)
	if err != nil {
		return err
	}
	c.K8sContainer.VolumeMounts = vm

	return nil
}

//prs, err := dbCont.parseProbes(dbCont.EnvVar)
//k8sCont = env

func (d *DbContainers) parseProbes(probeString string) (*v1.Probe, error) {
	pr := &v1.Probe{}
	err := json.Unmarshal([]byte(probeString), &pr)
	if err != nil {
		return pr, err
	}
	return pr, err

}
func (d *DbContainers) parseContainerPorts(portName, portProtocol string, portNumber int32) []v1.ContainerPort {
	contPort := v1.ContainerPort{
		Name:          portName,
		HostPort:      0,
		ContainerPort: portNumber,
		Protocol:      v1.Protocol(portProtocol),
		HostIP:        "",
	}

	return []v1.ContainerPort{contPort}
}

func (d *DbContainers) parseVolumeMount(volName, volPathString string) ([]v1.VolumeMount, error) {
	var vmSingleK8s v1.VolumeMount

	vmSingleK8s.Name = volName
	vmSingleK8s.MountPath = volPathString

	return []v1.VolumeMount{vmSingleK8s}, nil
}

func (d *DbContainers) parseEnvVars(envVarString string) ([]v1.EnvVar, error) {
	var envVars []v1.EnvVar
	m := map[string]string{}
	err := json.Unmarshal([]byte(envVarString), &m)
	if err != nil {
		return envVars, err
	}
	for k, v := range m {
		envSource := v1.EnvVarSource{}
		verr := json.Unmarshal([]byte(v), &envSource)
		if verr != nil {
			return envVars, verr
		}
		envVar := v1.EnvVar{
			Name:      k,
			Value:     "",
			ValueFrom: &envSource,
		}
		envVars = append(envVars, envVar)
	}
	return envVars, err
}

type Containers struct {
	Containers []Container

	K8sContainerSlice []v1.Container
	DBContainerSlice  DBContainerSlice
}

type Container struct {
	containers.Container
	K8sContainer v1.Container
	DB           DbContainers
}

type DBContainerSlice []DbContainers
type DbContainers struct {
	PortName     string
	PortProtocol string
	PortNumber   int
	EnvVar       string
	Probes       string
	VolumeName   string
	VolumePath   string
}

//func(c *Containers) ParseDBFields() apps.RowValues {
//
//	for _, dbCont := range c.DBContainerSlice {
//		dbCont.parseContainerPorts()
//	}
//
//
//	return
//}

func NewContainer() Container {
	return Container{
		containers.NewContainer(), NewK8sContainer(), NewDbContainers(),
	}
}

func NewK8sContainer() v1.Container {
	return v1.Container{}
}
func NewDbContainers() DbContainers {
	return DbContainers{}
}
