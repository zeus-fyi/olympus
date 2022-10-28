package ingresses

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

type IngressClassName struct {
	structs.ChildClassSingleValue
}

func (is *Spec) NewIngressClassName(ingressClassName string) {
	csv := structs.NewChildClassSingleValue("ingressClassName")
	ts := chronos.Chronos{}
	csv.SetSingleChildClassIDTypeNameKeyAndValue(ts.UnixTimeStampNow(), "ingressClassName", "ingressClassName", ingressClassName)

	icName := IngressClassName{csv}
	is.IngressClassName = &icName
}
