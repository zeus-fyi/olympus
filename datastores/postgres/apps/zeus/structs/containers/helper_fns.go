package containers

import (
	"fmt"
	"strings"

	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

func (c *Container) ProcessAndSetAmbiguousContainerFieldStatusAndSubfieldIds() {
	timeStamper := chronos.Chronos{}
	if len(c.Metadata.ContainerPlatformOs) <= 0 {
		c.Metadata.ContainerPlatformOs = fmt.Sprintf("%d-latest", timeStamper.UnixTimeStampNow())
	}
	if len(c.Metadata.ContainerRepository) <= 0 {
		c.Metadata.ContainerRepository = "undefined"
	}
	if len(c.Metadata.ContainerVersionTag) <= 0 {
		c.Metadata.ContainerVersionTag = fmt.Sprintf("%d-latest", timeStamper.UnixTimeStampNow())
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
	for i, _ := range c.Ports {
		c.Ports[i].PortID = timeStamper.UnixTimeStampNow()
	}
	for i, _ := range c.Env {
		c.Env[i].EnvID = timeStamper.UnixTimeStampNow()
	}
	for i, _ := range c.VolumeMounts {
		c.VolumeMounts[i].VolumeMountID = timeStamper.UnixTimeStampNow()
	}
	for i, _ := range c.Probes {
		c.Probes[i].ProbeID = timeStamper.UnixTimeStampNow()
	}
	if c.ResourceRequest != nil {
		c.ResourceRequest.ComputeResourcesID = timeStamper.UnixTimeStampNow()
	}
}
