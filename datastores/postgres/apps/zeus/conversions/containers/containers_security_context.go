package containers

import (
	"encoding/json"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/containers"
	v1 "k8s.io/api/core/v1"
)

func ConvertContainerSecurityContextToContainerDB(c v1.Container, dbContainer containers.Container) (containers.Container, error) {
	securityContext := c.SecurityContext

	if securityContext != nil {
		b, err := json.Marshal(securityContext)
		if err != nil {
			return dbContainer, err
		}
		dbContainer.SecurityContext.SecurityContextKeyValues = string(b)
	}
	return dbContainer, nil
}
