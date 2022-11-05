package temporal_client

import (
	"crypto/tls"

	"go.temporal.io/sdk/client"
)

func ConnectClient(clientCertPath, clientKeyPath string) error {
	cert, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
	if err != nil {
		return err
	}
	temporalClient, err := client.Dial(
		client.Options{
			HostPort:  "your-custom-namespace.tmprl.cloud:7233",
			Namespace: "your-custom-namespace",
			ConnectionOptions: client.ConnectionOptions{
				TLS: &tls.Config{Certificates: []tls.Certificate{cert}},
			},
		})
	defer temporalClient.Close()
	return err
}
