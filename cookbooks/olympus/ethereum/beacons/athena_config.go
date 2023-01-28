package olympus_beacon_cookbooks

import (
	zeus_topology_config_drivers "github.com/zeus-fyi/zeus/pkg/zeus/workload_config_drivers"
	v1 "k8s.io/api/core/v1"
)

// TODO env vars

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
			Command:         []string{"/bin/sh"},
			Args:            AthenaCmdArgs(),
			Ports:           []v1.ContainerPort{AthenaPort},
			EnvFrom:         nil,
			Env:             nil,
			Resources:       v1.ResourceRequirements{},
			VolumeMounts:    nil,
			VolumeDevices:   nil,
			LivenessProbe:   nil,
			ReadinessProbe:  nil,
			StartupProbe:    nil,
			ImagePullPolicy: "Always",
			SecurityContext: nil,
			Stdin:           false,
			StdinOnce:       false,
			TTY:             false,
		}}
)

const (
	startDashC     = "-c"
	athenaStartCmd = "athena"
)

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
