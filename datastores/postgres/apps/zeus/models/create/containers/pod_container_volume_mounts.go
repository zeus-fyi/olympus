package containers

import (
	"fmt"

	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/containers"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func CreateVolumeMountsCTEs() (sql_query_templates.SubCTE, sql_query_templates.SubCTE) {
	// env vars
	ts := chronos.Chronos{}
	vmContainer := autogen_bases.ContainerVolumeMounts{}
	contVmsSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_containers_volume_mounts_%d", ts.UnixTimeStampNow()))
	contVmsSubCTE.TableName = vmContainer.GetTableName()
	contVmsSubCTE.Columns = []string{"volume_mount_id", "volume_mount_path", "volume_name", "volume_read_only", "volume_sub_path"}
	contVmsRelationshipsSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_containers_volume_mounts_relationships_%d", ts.UnixTimeStampNow()))
	contVmsRelationshipsSubCTE.TableName = "containers_volume_mounts"
	contVmsRelationshipsSubCTE.Columns = []string{"chart_subcomponent_child_class_type_id", "container_id", "volume_mount_id"}
	return contVmsSubCTE, contVmsRelationshipsSubCTE
}

func (p *PodTemplateSpec) insertContainerVolumeMountsValues(m map[string]containers.Container, imageID string, contVmsSubCTE, contVmsRelationshipsSubCTE *sql_query_templates.SubCTE) {
	c, ok := m[imageID]
	if !ok {
		return
	}
	ts := chronos.Chronos{}
	podSpecChildClassTypeID := p.GetPodSpecChildClassTypeID()

	for _, vm := range c.VolumeMounts {
		vmID := ts.UnixTimeStampNow()
		vm.VolumeMountID = vmID
		contVmsSubCTE.AddValues(vmID, vm.VolumeMountPath, vm.VolumeName, vm.VolumeReadOnly, vm.VolumeSubPath)
		contVmsRelationshipsSubCTE.AddValues(podSpecChildClassTypeID, c.GetContainerID(), vmID)
	}
	return
}
