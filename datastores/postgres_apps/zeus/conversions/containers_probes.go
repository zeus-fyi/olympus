package conversions

import (
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/autogen"
	v1 "k8s.io/api/core/v1"
)

func ContainerProbesToDB(cs *v1.Container) []autogen_structs.ContainerProbes {
	suProbe := cs.StartupProbe
	livenessProbe := cs.LivenessProbe
	readinessProbe := cs.ReadinessProbe
	// from k8s
	probes := probesThatExist(suProbe, livenessProbe, readinessProbe)
	// to db format
	probeSlice := make([]autogen_structs.ContainerProbes, len(probes))

	for i, p := range probes {
		probe := ContainerProbeToDB(p)
		probeSlice[i] = probe
	}
	return probeSlice
}

func ContainerProbeToDB(p *v1.Probe) autogen_structs.ContainerProbes {
	dbProbe := autogen_structs.ContainerProbes{
		ProbeID:             0,
		ProbeKeyValuesJSONb: "",
	}
	return dbProbe
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
