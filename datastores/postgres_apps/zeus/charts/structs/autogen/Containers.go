package autogen_structs

type Containers struct {
	ContainerID              int    `db:"container_id"`
	ContainerName            string `db:"container_name"`
	ContainerImageID         string `db:"container_image_id"`
	ContainerVersionTag      string `db:"container_version_tag"`
	ContainerPlatformOs      string `db:"container_platform_os"`
	ContainerRepository      string `db:"container_repository"`
	ContainerImagePullPolicy string `db:"container_image_pull_policy"`
}
