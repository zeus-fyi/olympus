package containers

func (p *PodContainersGroup) insertContainerToPodSpec() string {
	return "INSERT INTO chart_subcomponent_spec_pod_template_containers(chart_subcomponent_child_class_type_id, container_id, is_init_container, container_sort_order) VALUES "
}
