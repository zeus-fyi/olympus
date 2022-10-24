package containers

import (
	"fmt"
)

func selectRelatedContainerIDFromImageID(imageID string) string {
	selectRelatedContainerID := fmt.Sprintf("(SELECT container_id FROM containers WHERE container_image_id = '%s' ORDER BY container_id DESC LIMIT 1)", imageID)
	return selectRelatedContainerID
}
