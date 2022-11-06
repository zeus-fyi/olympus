package temporal_client

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

func ConnectClient(authCfg TemporalAuth) error {
	cert, err := tls.LoadX509KeyPair(authCfg.ClientCertPath, authCfg.ClientPEMKeyPath)
	if err != nil {
		return err
	}
	temporalClient, err := client.Dial(
		client.Options{
			HostPort:  authCfg.HostPort,
			Namespace: authCfg.Namespace,
			ConnectionOptions: client.ConnectionOptions{
				TLS: &tls.Config{Certificates: []tls.Certificate{cert}},
			},
		})
	defer temporalClient.Close()
	return err
}
