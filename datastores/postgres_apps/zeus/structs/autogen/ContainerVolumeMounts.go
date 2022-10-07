package autogen_structs

type ContainerVolumeMounts struct {
	VolumeMountID   int    `db:"volume_mount_id"`
	VolumeMountPath string `db:"volume_mount_path"`
	VolumeName      string `db:"volume_name"`
}
