package common

import (
	"fmt"
	"strings"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/containers"
)

// TODO depcreate
func InsertContainerValues(parentExpression string, containers containers.Containers) string {
	valsToInsert := "VALUES "
	for i, cont := range containers {
		c := cont.Metadata
		// todo ports, env, probes

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

		valsToInsert += fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', '%s')", c.ContainerName, c.ContainerImageID, c.ContainerVersionTag, c.ContainerPlatformOs, c.ContainerRepository, c.ContainerImagePullPolicy)
		if i < len(containers)-1 {
			valsToInsert += ","
		}
	}

	containerInsert := fmt.Sprintf(`
				%s AS (
					INSERT INTO containers(container_name, container_image_id, container_version_tag, container_platform_os, container_repository, container_image_pull_policy)
					%s
	),`, "cte_containers", valsToInsert)

	returnExpression := fmt.Sprintf("%s %s", parentExpression, containerInsert)
	return returnExpression
}
