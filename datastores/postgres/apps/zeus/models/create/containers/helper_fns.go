package containers

import (
	"fmt"
	"strings"

	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
)

func processAndSetAmbiguousContainerFieldStatus(c autogen_bases.Containers) {
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

func selectRelatedContainerIDFromImageID(imageID string) string {
	selectRelatedContainerID := fmt.Sprintf("SELECT container_id FROM containers WHERE container_image_id = '%s'", imageID)
	return selectRelatedContainerID
}

func AppendValuesListComma(parentExpression string, valuesInsertCommaBoolCond bool) string {
	if valuesInsertCommaBoolCond {
		parentExpression += ","
	}
	return parentExpression
}

func (p *PodContainersGroup) process() {

}
