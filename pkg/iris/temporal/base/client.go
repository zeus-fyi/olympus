package temporal_base

import (
	"crypto/tls"

	"go.temporal.io/sdk/client"
)

type TemporalAuth struct {
	ClientCertPath   string
	ClientPEMKeyPath string
	Namespace        string
	HostPort         string
}

type TemporalClient struct {
	client.Client
	client.Options
}

// 	ConnectClient must defer temporalClient.Close()
func NewTemporalClient(authCfg TemporalAuth) (TemporalClient, error) {
	tc := TemporalClient{}
	cert, err := tls.LoadX509KeyPair(authCfg.ClientCertPath, authCfg.ClientPEMKeyPath)
	if err != nil {
		return tc, err
	}

	opts := client.Options{
		HostPort:  authCfg.HostPort,
		Namespace: authCfg.Namespace,
		ConnectionOptions: client.ConnectionOptions{
			TLS: &tls.Config{Certificates: []tls.Certificate{cert}},
		},
	}
	tc.Options = opts
	return tc, err
}
