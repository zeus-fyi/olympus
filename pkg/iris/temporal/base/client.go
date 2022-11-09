package temporal_base

import (
	"crypto/tls"

	"github.com/rs/zerolog"
	"github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	"go.temporal.io/sdk/client"
	zerologadapter "logur.dev/adapter/zerolog"
	"logur.dev/logur"
)

type TemporalClient struct {
	client.Client
	client.Options
}

// NewTemporalClient must call to connect and then must defer temporalClient.Close()
func NewTemporalClient(authCfg auth.TemporalAuth) (TemporalClient, error) {
	tc := TemporalClient{}
	cert, err := tls.LoadX509KeyPair(authCfg.ClientCertPath, authCfg.ClientPEMKeyPath)
	if err != nil {
		return tc, err
	}
	logger := logur.LoggerToKV(zerologadapter.New(zerolog.Nop()))

	opts := client.Options{
		Logger:    logger,
		HostPort:  authCfg.HostPort,
		Namespace: authCfg.Namespace,
		ConnectionOptions: client.ConnectionOptions{
			TLS: &tls.Config{Certificates: []tls.Certificate{cert}},
		},
	}
	tc.Options = opts
	return tc, err
}
