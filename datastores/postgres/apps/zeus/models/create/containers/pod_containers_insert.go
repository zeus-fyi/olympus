package containers

import "fmt"

// insertPodContainerGroupSQL, will use the next_id distributed ID generator and select the container id
// value for subsequent subcomponent relationships of its element, should greatly simplify the insert logic
func (p *PodContainersGroup) insertPodContainerGroupSQL() string {
	q := fmt.Sprintf(
		`WITH cte_insert_containers AS (
					%s
				), cte_insert_spec AS (
					%s
				), `,
	)

	return q
}

func selectRelatedContainerIDFromImageID(imageID string) string {
	selectRelatedContainerID := fmt.Sprintf("SELECT container_id FROM containers WHERE container_image_id = %s", imageID)
	return selectRelatedContainerID
}
