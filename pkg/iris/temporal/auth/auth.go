package temporal_auth

type TemporalAuth struct {
	ClientCertPath   string
	ClientPEMKeyPath string
	ServerRootCACert string
	Namespace        string
	HostPort         string
}
