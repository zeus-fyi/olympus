package create

import "fmt"

func (d *Deployment) insertDeploymentMetadataChildren() string {

	//m := d.Metadata
	//if len(m.Name) > 0 {
	//
	//}
	//
	//if len(m.Labels) > 0 {
	//
	//}
	//
	//if len(m.Annotations) > 0 {
	//
	//}
	s := fmt.Sprintf(
		`cte_insert_metadata_children AS (
					INSERT INTO chart_subcomponent_child_class_types(chart_subcomponent_parent_class_type_id, chart_subcomponent_child_class_type_name)
					VALUES ((SELECT chart_subcomponent_parent_class_type_id FROM cte_insert_metadata), '%s')
					RETURNING chart_subcomponent_child_class_type_id
	)`, "c")
	return s
}

func (d *Deployment) insertChildClassType() string {
	s := fmt.Sprintf(
		`cte_insert_metadata_children AS (
					INSERT INTO chart_subcomponent_child_class_types(chart_subcomponent_parent_class_type_id, chart_subcomponent_child_class_type_name)
					VALUES ((SELECT chart_subcomponent_parent_class_type_id FROM cte_insert_metadata), '%s')
					RETURNING chart_subcomponent_child_class_type_id
	)`, "c")
	return s
}
