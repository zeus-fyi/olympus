package containers

type DBContainerSlice []DbContainers
type DbContainers struct {
	Ports            string
	EnvVar           string
	Probes           string
	ContainerVolumes string
}
