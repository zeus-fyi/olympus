package networking

import (
	v1 "k8s.io/api/core/v1"

	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/conversions/common"
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/networking"
)

func ConvertServiceConfigToDB(svc *v1.Service) networking.Service {
	dbService := networking.NewService()
	dbService.Metadata = common.CreateMetadataByFields(svc.Name, svc.Annotations, svc.Labels)
	return dbService
}
