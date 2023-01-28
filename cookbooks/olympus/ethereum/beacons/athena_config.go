package olympus_beacon_cookbooks

import (
	zeus_topology_config_drivers "github.com/zeus-fyi/zeus/pkg/zeus/workload_config_drivers"
	v1 "k8s.io/api/core/v1"
)

var (
	AthenaContainer = zeus_topology_config_drivers.ContainerDriver{
		Container: v1.Container{
			Name:                     "",
			Image:                    "",
			Command:                  nil,
			Args:                     nil,
			WorkingDir:               "",
			Ports:                    nil,
			EnvFrom:                  nil,
			Env:                      nil,
			Resources:                v1.ResourceRequirements{},
			VolumeMounts:             nil,
			VolumeDevices:            nil,
			LivenessProbe:            nil,
			ReadinessProbe:           nil,
			StartupProbe:             nil,
			Lifecycle:                nil,
			TerminationMessagePath:   "",
			TerminationMessagePolicy: "",
			ImagePullPolicy:          "Always",
			SecurityContext:          nil,
			Stdin:                    false,
			StdinOnce:                false,
			TTY:                      false,
		}}
)
