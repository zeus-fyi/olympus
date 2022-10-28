package create_ingresses

import "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking/ingresses"

type Ingress struct {
	ingresses.Ingress
}

func NewCreateIngress() Ingress {
	return Ingress{ingresses.NewIngress()}
}

func (i *Ingress) InsertIngress() {
	// TODO
}
