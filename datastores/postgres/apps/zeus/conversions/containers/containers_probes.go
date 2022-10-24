package containers

import (
	"encoding/json"

	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/containers"
	v1 "k8s.io/api/core/v1"
)

func ContainerProbeToDB(p *v1.Probe) (autogen_bases.ContainerProbes, error) {
	dbProbe := autogen_bases.ContainerProbes{
		ProbeID:             0,
		ProbeKeyValuesJSONb: "",
	}
	if p != nil {
		bytes, err := json.Marshal(p)
		if err != nil {
			return dbProbe, err
		}
		dbProbe.ProbeKeyValuesJSONb = string(bytes)
	}

	return dbProbe, nil
}

func probesThatExist(probes ...*v1.Probe) []*v1.Probe {
	var probeSlice []*v1.Probe
	for _, probe := range probes {
		if probe != nil {
			probeSlice = append(probeSlice, probe)
		}
	}
	return probeSlice
}

func ConvertContainerProbesToDB(cs v1.Container, dbContainer containers.Container) (containers.Container, error) {
	suProbe := cs.StartupProbe
	livenessProbe := cs.LivenessProbe
	readinessProbe := cs.ReadinessProbe
	// from k8s
	probes := probesThatExist(suProbe, livenessProbe, readinessProbe)
	// to db format
	probeSlice := make(autogen_bases.ContainerProbesSlice, len(probes))

	for i, p := range probes {
		probe, err := ContainerProbeToDB(p)
		if err != nil {
			return dbContainer, err
		}
		probeSlice[i] = probe
	}
	dbContainer.Probes = probeSlice
	return dbContainer, nil
}
