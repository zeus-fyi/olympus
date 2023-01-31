package containers

import (
	"encoding/json"
	"strings"

	"github.com/rs/zerolog/log"
	v1 "k8s.io/api/core/v1"
)

func (c *Container) ParseFields() error {
	var k8sCont v1.Container

	err := c.DB.parseCmdArgs(&k8sCont, c.DB.CmdArgs)
	if err != nil {
		return err
	}
	err = c.DB.parseSecurityContext(&k8sCont, c.DB.SecurityContext)
	if err != nil {
		return err
	}
	err = c.DB.parseComputeResources(&k8sCont, c.DB.ComputeResources)
	if err != nil {
		return err
	}
	err = c.DB.parseProbes(&k8sCont, c.DB.Probes)
	if err != nil {
		return err
	}
	env, err := c.DB.parseEnvVars(c.DB.EnvVar)
	if err != nil {
		return err
	}
	k8sCont.Env = env

	ports, err := c.DB.parseContainerPorts(c.DB.Ports)
	if err != nil {
		return err
	}
	k8sCont.Ports = ports
	contVms, err := c.DB.parseVolumeMount(c.DB.ContainerVolumes)
	if err != nil {
		return err
	}

	k8sCont.VolumeMounts = contVms
	c.K8sContainer = k8sCont
	return nil
}

func (d *DbContainers) parseSecurityContext(container *v1.Container, securityContext string) error {
	m := make(map[string]string)
	err := json.Unmarshal([]byte(securityContext), &m)
	if err != nil {
		log.Info().Interface("securityContext", securityContext)
		log.Err(err).Msg("DbContainers: parseSecurityContext unmarshal securityContext string")
		return err
	}
	for _, nv := range m {
		if len(nv) > 0 {
			perr := json.Unmarshal([]byte(nv), &container.SecurityContext)
			if perr != nil {
				log.Err(err).Msg("DbContainers: parseSecurityContext Unmarshal &container.SecurityContext")
				return err
			}
		}
	}
	return err
}

func (d *DbContainers) parseCmdArgs(container *v1.Container, cmdArgs string) error {
	m := make(map[string]map[string]interface{})
	err := json.Unmarshal([]byte(cmdArgs), &m)
	if err != nil {
		log.Err(err).Msg("DbContainers: parseCmdArgs unmarshal cmdArgs string")
		return err
	}
	for _, v := range m {
		for nk, nv := range v {
			if len(nv.(string)) > 0 {
				switch nk {
				case "command":
					container.Command = strings.Split(nv.(string), ",")
				case "args":
					container.Args = strings.Split(nv.(string), ",")
				}
			}
		}
	}
	return err
}

func (d *DbContainers) parseComputeResources(container *v1.Container, computeResources string) error {
	m := make(map[string]map[string]interface{})
	err := json.Unmarshal([]byte(computeResources), &m)
	if err != nil {
		log.Info().Interface("computeResources", computeResources)
		log.Err(err).Msg("DbContainers: parseComputeResources unmarshal parseComputeResources string")
		return err
	}
	for _, v := range m {
		for nk, nv := range v {
			bytes, berr := json.Marshal(nv)
			if berr != nil {
				log.Err(err).Msg("DbContainers: parseComputeResources Marshal")
				return berr
			}
			rl := v1.ResourceList{}
			perr := json.Unmarshal(bytes, &rl)
			if perr != nil {
				log.Err(err).Msg("DbContainers: parseComputeResources Unmarshal to ResourceList")
				return err
			}
			if rl.Cpu().Value() == int64(0) {
				delete(rl, "cpu")
			}
			if rl.Memory().Value() == int64(0) {
				delete(rl, "memory")
			}
			if rl.StorageEphemeral().Value() == int64(0) {
				delete(rl, "ephemeral-storage")
			}
			switch nk {
			case "limits":
				container.Resources.Limits = rl
			case "requests":
				container.Resources.Requests = rl
			}
		}
	}
	return err
}

func (d *DbContainers) parseProbes(container *v1.Container, probeString string) error {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(probeString), &m)
	if err != nil {
		log.Info().Interface("probeString", probeString)
		log.Err(err).Msg("DbContainers: parseProbes Unmarshal probe string")
		return err
	}
	for k, v := range m {
		pr := &v1.Probe{}
		bytes, berr := json.Marshal(v)
		if berr != nil {
			log.Info().Interface("v", v)
			log.Err(err).Msg("DbContainers: parseProbes Marshal")
			return berr
		}
		perr := json.Unmarshal(bytes, &pr)
		if perr != nil {
			log.Info().Interface("prBytes", bytes)
			log.Err(err).Msg("DbContainers: parseProbes Unmarshal")
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
		log.Info().Interface("portsStr", portsStr)
		log.Err(err).Msg("DbContainers: parseContainerPorts")
		return ports, err
	}
	for _, v := range m {
		bytes, berr := json.Marshal(v)
		if berr != nil {
			log.Info().Interface("v", v)
			log.Err(err).Msg("DbContainers: parseContainerPorts")
			return ports, berr
		}
		var port v1.ContainerPort
		perr := json.Unmarshal(bytes, &port)
		if perr != nil {
			log.Info().Interface("portBytes", bytes)
			log.Err(err).Msg("DbContainers: parseContainerPorts")
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
		log.Info().Interface("contVolMounts", contVolMounts)
		log.Err(err).Msg("DbContainers: parseVolumeMount")
		return contVms, err
	}
	for _, v := range m {
		bytes, berr := json.Marshal(v)
		if berr != nil {
			log.Info().Interface("v", v)
			log.Err(err).Msg("DbContainers: parseVolumeMount Marshal")
			return contVms, berr
		}
		var contVm v1.VolumeMount
		perr := json.Unmarshal(bytes, &contVm)
		if perr != nil {
			log.Info().Interface("contVm", contVm)
			log.Err(err).Msg("DbContainers: parseVolumeMount Unmarshal to VolumeMount")
			return contVms, perr
		}
		contVms = append(contVms, contVm)
	}
	return contVms, nil
}

func (d *DbContainers) parseEnvVars(envVarString string) ([]v1.EnvVar, error) {
	var envVars []v1.EnvVar
	m := map[string]any{}
	err := json.Unmarshal([]byte(envVarString), &m)
	if err != nil {
		log.Info().Interface("envVarString", envVarString)
		log.Err(err).Msg("DbContainers: parseEnvVars")
		return envVars, err
	}
	for k, v := range m {
		envSource := v1.EnvVarSource{}
		verr := json.Unmarshal([]byte(v.(string)), &envSource)
		if verr != nil {
			log.Err(err).Msg("DbContainers: parseEnvVars")
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
