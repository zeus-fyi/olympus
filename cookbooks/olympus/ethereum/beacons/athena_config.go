package olympus_beacon_cookbooks

import (
	"fmt"

	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	zeus_topology_config_drivers "github.com/zeus-fyi/zeus/zeus/workload_config_drivers/config_overrides"
	v1 "k8s.io/api/core/v1"
)

const (
	protocolNetworkKeyEnv = "PROTOCOL_NETWORK_ID"
	startDashC            = "-c"
	athenaStartCmd        = "athena"
)

var (
	AthenaPort = v1.ContainerPort{
		Name:          "athena",
		ContainerPort: 9003,
		Protocol:      v1.Protocol("TCP"),
	}
	AthenaContainer = zeus_topology_config_drivers.ContainerDriver{
		Container: v1.Container{
			Name:            "athena",
			Image:           "registry.digitalocean.com/zeus-fyi/athena:latest",
			Ports:           []v1.ContainerPort{AthenaPort},
			ImagePullPolicy: "Always",
		}}
)

func AthenaContainerDriver(network string) zeus_topology_config_drivers.ContainerDriver {
	switch network {
	case "mainnet":
		envVar := AthenaContainer.CreateEnvVarKeyValue(protocolNetworkKeyEnv, fmt.Sprintf("%d", hestia_req_types.EthereumMainnetProtocolNetworkID))
		AthenaContainer.AppendEnvVars = []v1.EnvVar{envVar}
	case "ephemery":
		envVar := AthenaContainer.CreateEnvVarKeyValue(protocolNetworkKeyEnv, fmt.Sprintf("%d", hestia_req_types.EthereumEphemeryProtocolNetworkID))
		AthenaContainer.AppendEnvVars = []v1.EnvVar{envVar}
	}
	return AthenaContainer
}

func AthenaCmdArgs() []string {
	args := []string{startDashC, athenaStartCmd}

	args = append(args, "--bearer=\"${BEARER}\"")
	args = append(args, "--cloud-provider=\"${CLOUD_PROVIDER}\"")
	args = append(args, "--ctx=\"${CTX}\"")
	args = append(args, "--ns=\"${NS}\"")
	args = append(args, "--region=\"${REGION}\"")

	args = append(args, "--age-private-key=\"${AGE_PKEY}\"")
	args = append(args, "--do-spaces-key=\"${DO_SPACES_KEY}\"")
	args = append(args, "--do-spaces-private-key=\"${DO_SPACES_PKEY}\"")
	args = append(args, "--env=\"production\"")

	return args
}
