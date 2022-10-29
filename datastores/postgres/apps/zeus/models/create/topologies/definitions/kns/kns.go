package create_kns

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
)

type Kns struct {
	kns.Kns
}

func NewCreateKns() Kns {
	ckns := Kns{kns.NewKns()}
	return ckns
}
