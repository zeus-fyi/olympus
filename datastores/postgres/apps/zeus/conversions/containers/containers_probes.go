package containers

import (
	"encoding/json"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/containers"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/containers/probes"
	v1 "k8s.io/api/core/v1"
)

func ContainerProbeToDB(p *v1.Probe) (probes.Probe, error) {
	dbProbe := probes.Probe{}
	if p != nil {
		bytes, err := json.Marshal(p)
		if err != nil {
			return dbProbe, err
		}
		dbProbeKeyValuesJsonb := string(bytes)
		dbProbe.ProbeKeyValuesJSONb = dbProbeKeyValuesJsonb
	}

	return dbProbe, nil
}

// expected order suProbe, livenessProbe, readinessProbe
func probesThatExist(prs ...*v1.Probe) (probes.ProbeSlice, error) {
	var probeSlice probes.ProbeSlice
	for i, probe := range prs {
		if probe != nil {
			p, err := ContainerProbeToDB(probe)
			if err != nil {
				return probeSlice, err
			}
			switch i {
			case 0:
				p.ProbeType = "startupProbe"
			case 1:
				p.ProbeType = "livenessProbe"
			case 2:
				p.ProbeType = "readinessProbe"
			}
			probeSlice = append(probeSlice, p)
		}
	}
	return probeSlice, nil
}

func ConvertContainerProbesToDB(cs v1.Container, dbContainer containers.Container) (containers.Container, error) {
	suProbe := cs.StartupProbe
	livenessProbe := cs.LivenessProbe
	readinessProbe := cs.ReadinessProbe
	// from k8s
	existingProbes, prErr := probesThatExist(suProbe, livenessProbe, readinessProbe)
	if prErr != nil {
		return dbContainer, prErr
	}
	// to db format
	dbContainer.Probes = existingProbes
	return dbContainer, nil
}
