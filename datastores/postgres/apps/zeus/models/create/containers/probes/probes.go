package probes

import autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"

type Probe struct {
	autogen_bases.ContainersProbes
	autogen_bases.ContainerProbes
}

type ProbeSlice []Probe

func (p *Probe) SetContainerID(id int) {
	p.ContainerID = id
}

func (p *Probe) SetProbeID(id int) {
	p.ContainerProbes.ProbeID = id
	p.ContainersProbes.ProbeID = id
}
