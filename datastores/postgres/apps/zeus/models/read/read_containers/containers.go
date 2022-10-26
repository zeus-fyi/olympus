package read_containers

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/containers"
	v1 "k8s.io/api/core/v1"
)

func (c *Container) ParseFields() error {
	var k8sCont v1.Container
	err := c.DB.parseProbes(&k8sCont, c.DB.Probes)
	if err != nil {
		return err
	}
	env, err := c.DB.parseEnvVars(c.DB.EnvVar)
	if err != nil {
		return err
	}
	k8sCont.Env = env

	ports, err := c.DB.parseContainerPorts(c.DB.Ports)
	if err != nil {
		return err
	}
	k8sCont.Ports = ports
	contVms, err := c.DB.parseVolumeMount(c.DB.ContainerVolumes)
	if err != nil {
		return err
	}
	k8sCont.VolumeMounts = contVms
	c.K8sContainer = k8sCont

	return nil
}

type Containers struct {
	Containers []Container

	DBContainerSlice DBContainerSlice
}

type Container struct {
	containers.Container
	DB DbContainers
}

func NewContainer() Container {
	return Container{
		containers.NewContainer(), NewDbContainers(),
	}
}

func NewK8sContainer() v1.Container {
	return v1.Container{}
}
func NewDbContainers() DbContainers {
	return DbContainers{}
}
