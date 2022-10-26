package read_containers

import (
	"encoding/json"

	v1 "k8s.io/api/core/v1"
)

func (d *DbContainers) parseProbes(container *v1.Container, probeString string) error {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(probeString), &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		pr := &v1.Probe{}

		bytes, berr := json.Marshal(v)
		if berr != nil {
			return berr
		}
		perr := json.Unmarshal(bytes, &pr)
		if perr != nil {
			return err
		}
		switch k {
		case "readinessProbe":
			container.ReadinessProbe = pr
		case "livenessProbe":
			container.LivenessProbe = pr
		case "startupProbe":
			container.StartupProbe = pr
		}
	}
	return err

}

// TODO this is joining on many, so need to remove duplicates
func (d *DbContainers) parseContainerPorts(portsStr string) ([]v1.ContainerPort, error) {
	m := make(map[string]interface{})
	var ports []v1.ContainerPort
	err := json.Unmarshal([]byte(portsStr), &m)
	if err != nil {
		return ports, err
	}
	for _, v := range m {
		bytes, berr := json.Marshal(v)
		if berr != nil {
			return ports, berr
		}
		var port v1.ContainerPort
		perr := json.Unmarshal(bytes, &port)
		if perr != nil {
			return ports, perr
		}
		ports = append(ports, port)
	}

	return ports, nil
}

func (d *DbContainers) parseVolumeMount(contVolMounts string) ([]v1.VolumeMount, error) {
	m := make(map[string]interface{})
	var contVms []v1.VolumeMount
	err := json.Unmarshal([]byte(contVolMounts), &m)
	if err != nil {
		return contVms, err
	}
	for _, v := range m {
		bytes, berr := json.Marshal(v)
		if berr != nil {
			return contVms, berr
		}
		var contVm v1.VolumeMount
		perr := json.Unmarshal(bytes, &contVm)
		if perr != nil {
			return contVms, perr
		}
		contVms = append(contVms, contVm)
	}
	return contVms, nil
}

func (d *DbContainers) parseEnvVars(envVarString string) ([]v1.EnvVar, error) {
	var envVars []v1.EnvVar
	m := map[string]string{}
	err := json.Unmarshal([]byte(envVarString), &m)
	if err != nil {
		return envVars, err
	}
	for k, v := range m {
		envSource := v1.EnvVarSource{}
		verr := json.Unmarshal([]byte(v), &envSource)
		if verr != nil {
			return envVars, verr
		}
		envVar := v1.EnvVar{
			Name:      k,
			Value:     "",
			ValueFrom: &envSource,
		}
		envVars = append(envVars, envVar)
	}
	return envVars, err
}
