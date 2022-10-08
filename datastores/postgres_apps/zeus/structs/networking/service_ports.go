package networking

type ServicePort struct {
	Name       string
	Protocol   string
	Port       int
	TargetPort int
	NodePort   int
}

type ServicePorts []ServicePort
