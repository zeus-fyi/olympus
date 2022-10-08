package networking

import (
	v1 "k8s.io/api/core/v1"

	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/conversions/common"
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/networking"
)

func ConvertServiceConfigToDB(svc *v1.Service) networking.Service {
	dbService := networking.NewService()
	dbService.Metadata = common.CreateMetadataByFields(svc.Name, svc.Annotations, svc.Labels)
	dbService.ServiceSpec = ConvertServiceSpecConfigToDB(svc)
	return dbService
}

func ConvertServiceSpecConfigToDB(svc *v1.Service) networking.ServiceSpec {
	dbServiceSpec := networking.ServiceSpec{
		Selector: common.ConvertSelectorByFields(svc.Spec.Selector),
		Ports:    ServicePortsToDB(svc.Spec.Ports),
	}
	return dbServiceSpec
}
