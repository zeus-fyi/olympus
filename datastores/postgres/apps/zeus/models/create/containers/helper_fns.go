package containers

import (
	"strings"

	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
)

func processAndSetAmbiguousContainerFieldStatus(c autogen_structs.Containers) {
	if len(c.ContainerPlatformOs) <= 0 {
		c.ContainerPlatformOs = "undefined"
	}
	if len(c.ContainerRepository) <= 0 {
		c.ContainerRepository = "undefined"
	}
	if len(c.ContainerVersionTag) <= 0 {
		c.ContainerVersionTag = "undefined"
	}

	splitValues := strings.Split(c.ContainerImageID, ":")

	if len(splitValues) == 2 {
		c.ContainerRepository = splitValues[0]
		c.ContainerVersionTag = splitValues[1]
	}
}
