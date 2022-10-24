package networking

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/networking"
	v1 "k8s.io/api/core/v1"
)

func ConvertServiceConfigToDB(svc *v1.Service) networking.Service {
	dbService := networking.NewService()
	dbService.Metadata = common_conversions.CreateMetadataByFields(svc.Name, svc.Annotations, svc.Labels)
	dbService.ServiceSpec = ConvertServiceSpecConfigToDB(svc)
	return dbService
}

func ConvertServiceSpecConfigToDB(svc *v1.Service) networking.ServiceSpec {
	dbServiceSpec := networking.ServiceSpec{
		Type:     common_conversions.ConvertKeyValueStringToChildValues("type", string(svc.Spec.Type)),
		Selector: common_conversions.ConvertSelectorByFields(svc.Spec.Selector),
		Ports:    ServicePortsToDB(svc.Spec.Ports),
	}
	return dbServiceSpec
}
