package temporal_auth

type TemporalAuth struct {
	ClientCertPath   string
	ClientPEMKeyPath string
	Namespace        string
	HostPort         string
	Bearer           string
}
