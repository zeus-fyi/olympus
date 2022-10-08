package networking

import (
	common2 "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common"
	networking2 "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/networking"
	v1 "k8s.io/api/core/v1"
)

func ConvertServiceConfigToDB(svc *v1.Service) networking2.Service {
	dbService := networking2.NewService()
	dbService.Metadata = common2.CreateMetadataByFields(svc.Name, svc.Annotations, svc.Labels)
	dbService.ServiceSpec = ConvertServiceSpecConfigToDB(svc)
	return dbService
}

func ConvertServiceSpecConfigToDB(svc *v1.Service) networking2.ServiceSpec {
	dbServiceSpec := networking2.ServiceSpec{
		Type:     common2.ConvertKeyValueStringToChildValues("type", string(svc.Spec.Type)),
		Selector: common2.ConvertSelectorByFields(svc.Spec.Selector),
		Ports:    ServicePortsToDB(svc.Spec.Ports),
	}
	return dbServiceSpec
}
