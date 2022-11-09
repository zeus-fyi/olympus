package temporal_base

import (
	"crypto/tls"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
	zerologadapter "logur.dev/adapter/zerolog"
	"logur.dev/logur"
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

// ConnectTemporalClient must have client conn closed when called
func (t *TemporalClient) ConnectTemporalClient() error {
	dial, err := client.Dial(t.Options)
	if err != nil {
		log.Err(err).Msg("ConnectTemporalClient: dial failed")
		return err
	}
	t.Client = dial
	return nil
}

// NewTemporalClient must call to connect and then must defer temporalClient.Close()
func NewTemporalClient(authCfg TemporalAuth) (TemporalClient, error) {
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
