package containers

type DBContainerSlice []DbContainers
type DbContainers struct {
	SecurityContext  string
	ComputeResources string
	CmdArgs          string
	Ports            string
	EnvVar           string
	Probes           string
	ContainerVolumes string
}
