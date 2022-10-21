package containers

import (
	"strings"

	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

func (c *Container) ProcessAndSetAmbiguousContainerFieldStatusAndSubfieldIds() {
	if len(c.Metadata.ContainerPlatformOs) <= 0 {
		c.Metadata.ContainerPlatformOs = "undefined"
	}
	if len(c.Metadata.ContainerRepository) <= 0 {
		c.Metadata.ContainerRepository = "undefined"
	}
	if len(c.Metadata.ContainerVersionTag) <= 0 {
		c.Metadata.ContainerVersionTag = "undefined"
	}
	splitValues := strings.Split(c.Metadata.ContainerImageID, ":")
	if len(splitValues) == 2 {
		c.Metadata.ContainerRepository = splitValues[0]
		c.Metadata.ContainerVersionTag = splitValues[1]
	}
	c.setContainerSubfieldIds()
}

func (c *Container) setContainerSubfieldIds() {
	timeStamper := chronos.Chronos{}
	commonIdTag := int(timeStamper.UnixTimeStampNow())
	for i, _ := range c.Ports {
		c.Ports[i].PortID = commonIdTag
	}

	for i, _ := range c.Env {
		c.Env[i].EnvID = commonIdTag
	}
	for i, _ := range c.VolumeMounts {
		c.VolumeMounts[i].VolumeMountID = commonIdTag
	}
	for i, _ := range c.Probes {
		c.Probes[i].ProbeID = commonIdTag
	}
	if c.ResourceRequest != nil {
		c.ResourceRequest.ComputeResourcesID = commonIdTag
	}
}
