package configuration

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

type Data struct {
	structs.SuperParentClass
}

func NewCMData(data map[string]string) Data {
	var ts chronos.Chronos
	pcID := ts.UnixTimeStampNow()
	childTypeID := ts.UnixTimeStampNow()

	cmData := structs.NewSuperParentClassWithMultiChildType("Data", "Data")
	cmData.SetParentClassTypeID(pcID)
	cmData.AddValues(data)
	cmData.SetChildClassTypeIDs(childTypeID)

	return Data{cmData}
}
