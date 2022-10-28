package ingress

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

func NewIngressClassName(ingressClassName string) structs.ChildClassSingleValue {
	csv := structs.NewChildClassSingleValue("ingressClassName")
	ts := chronos.Chronos{}
	csv.SetSingleChildClassIDTypeNameKeyAndValue(ts.UnixTimeStampNow(), "ingressClassName", "ingressClassName", ingressClassName)
	return csv
}
