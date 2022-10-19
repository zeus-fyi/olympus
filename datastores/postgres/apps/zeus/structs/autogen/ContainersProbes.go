package autogen_structs

type ContainersProbes struct {
	ProbeID     int    `db:"probe_id"`
	ContainerID int    `db:"container_id"`
	ProbeType   string `db:"probe_type"`
}
